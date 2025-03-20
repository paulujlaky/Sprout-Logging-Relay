import $ from 'jquery';

import { HideDialog, ShowContextMenu, ShowDialog, WaitForDialogResponse, type ContextMenuOption } from './Rendering';

const RelevantElements = {

    Document: $(document),

    Buttons: {

        AccountMenu: $(".DashHeaderAccount"),

    },

    Inputs: {

        Toggles: $("Toggle"),

    },

    Dialogs: {

    },

    Misc: {

        Branding: $(".DashHeaderBranding"),

    }
    
}

// General Caller

export function WatchPageInteractions(): void {

    HandleToggleClicks();

    HandleContextMenuInteractions();

    WatchBrandingClick();

    WatchForAccountMenu();

}

// Account Menu

function WatchForAccountMenu(): void {

    RelevantElements.Buttons.AccountMenu.on("click", async () => {

        const Dialog = $(".DashViewAccountDialog") as JQuery<HTMLDialogElement>;

        ShowDialog(Dialog);

        await WaitForDialogResponse(Dialog);

        HideDialog(Dialog);

    });

}

// Context Menu

const ContextMenuRouter = async (Event: JQuery.MouseEventBase | JQuery.TouchEventBase | JQuery.ContextMenuEvent): Promise<void> => {

    if (Event.type == "contextmenu") { Event.preventDefault(); } // Prevent default context menu

    // TODO: will need to determine target and if apt to show context menu

    return; // early return for now
    
    const SharedOptions: ContextMenuOption[] = []

    ShowContextMenu(Event, SharedOptions);

}

function HandleContextMenuInteractions(): void {

    RelevantElements.Document.on("contextmenu", ContextMenuRouter);
    RelevantElements.Document.on("click", ContextMenuRouter);

}

// Misc

function HandleToggleClicks(): void {

    RelevantElements.Inputs.Toggles.on("click", (Event) => {

        const Target = $(Event.target);

        const Toggle = Target.closest("Toggle");

        Toggle.toggleClass("Active");

    });

}

// Misc / Util

export function RemoveDialogOnEscape(Dialog: JQuery<HTMLDialogElement>): void {

    RelevantElements.Document.one("keydown", (Event) => {

        if (Event.key === "Escape") {

            HideDialog(Dialog);

        }

    });

}

export function SubmitDialogOnEnter(Dialog: JQuery<HTMLDialogElement>, EnterButton: JQuery<HTMLButtonElement>): void {

    const Listener = (Event: JQuery.KeyDownEvent) => {

        if (Event.key === "Enter") {

            EnterButton.trigger("click");

        }

    };

    RelevantElements.Document.on("keydown", Listener);

    Dialog.on("Close", () => {

        RelevantElements.Document.off("keydown", Listener);

    });
    
}

function WatchBrandingClick(): void {
    
    RelevantElements.Misc.Branding.on("click", () => {

        // Code here

    });

}
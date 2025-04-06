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

    Selects: {

        ServiceSelect: $("#ServiceSelect")

    },

    Misc: {

        Branding: $(".DashHeaderBranding"),

    }
    
}

// General Caller

export function WatchPageInteractions(): void {

    HandleToggleClicks();
    HandleSelectMenuInteractions();

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

function HandleSelectMenuInteractions(): void {

    const SelectMenus = $("SelectMenu");

    SelectMenus.on("click", (Event: JQuery.MouseEventBase | JQuery.TouchEventBase) => {

        const RelevantSelectMenu = $(Event.target).closest("SelectMenu");

        const SelectMenuIcon = RelevantSelectMenu.find(".SelectMenuIcon");
        const SelectMenuOptions = RelevantSelectMenu.find("SelectMenuOptions");

        const IsOpen = SelectMenuOptions.is(":visible");

        // Opens above the select menu

        const Coords: { X: number, Y: number } = {

            X: RelevantSelectMenu.position().left, // left of the select menu
            Y: ((RelevantSelectMenu.position().top - (SelectMenuOptions.outerHeight() || 0)) - 10) // top of the select menu - height of the options - 10px for spacing

        };

        SelectMenuOptions.css({ top: Coords.Y, left: Coords.X }); // set the position of the menu

        if (IsOpen) {

            SelectMenuOptions.hide();
            SelectMenuIcon.html("<ion-icon name='chevron-up'></ion-icon>");

        } else {

            SelectMenuOptions.show();
            SelectMenuIcon.html("<ion-icon name='chevron-down'></ion-icon>");

        }

    });        

}

interface SelectMenuOption {

    Name: string;
    Value: string;

}

export function WaitForSelectMenuResponse(SelectMenu: JQuery<HTMLElement>, OptionClickCallback: (Option: SelectMenuOption) => void): void {

    const Options = SelectMenu.find("SelectMenuOption"); // gets all INDIVIDUAL options

    const UpdateSelectMenuWithOption = (Option: SelectMenuOption): void => {

        SelectMenu.find(".SelectMenuLabel").text(Option.Name); // update the label with the new option name
        SelectMenu.attr("Value", Option.Value); // update the value of the select menu

    }

    Options.on("click", (Event: JQuery.MouseEventBase | JQuery.TouchEventBase) => {

        const RelatedOption = $(Event.target).closest("SelectMenuOption");

        const OptionValue = RelatedOption.attr("Value") || ""; // get the value of the option
        const OptionName = RelatedOption.text(); // the element text is the name

        const Option: SelectMenuOption = { Name: OptionName, Value: OptionValue };

        OptionClickCallback(Option); // call the callback with the option
        UpdateSelectMenuWithOption(Option); // update the select menu with the new option

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
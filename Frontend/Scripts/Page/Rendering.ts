import $ from "jquery";

import { AnimationTimes, type SproutAccount } from "../Misc/Structs";

import { RemoveDialogOnEscape, SubmitDialogOnEnter } from "./Interactions";

const RelevantElements = {

    Loading: {

        Container: $(".Container.DashLoading"),

    },

    Auth: {

        UsernameEmbeds: $(".DashHeaderMessageUser"),
        EmailEmbeds: $(".DashHeaderMessageEmail"),

    },

    Dialogs: { 
        
        Overlay: $(".DashDialogOverlay")

    },

    CtxMenu: $(".DashContextMenu")

}

// Loading

export function HideLoadingView(): void {

    RelevantElements.Loading.Container.fadeOut(AnimationTimes.Short);

}

// Auth

export function UpdateAuthElements(Account: SproutAccount) {

    RelevantElements.Auth.UsernameEmbeds.text(Account.Username);
    RelevantElements.Auth.EmailEmbeds.text(Account.Email);
    
}

// Dialogs

export function ShowDialog(Dialog: JQuery<HTMLDialogElement>): void {

    RelevantElements.Dialogs.Overlay.show();
    Dialog.css("display", "flex");

    // Clear inputs

    Dialog.find("input").val("").trigger("focus");
    Dialog.find("Toggle").removeClass("Active");

    const CloseBtn = Dialog.find(".DashDialogButtonCancel");
    const EnterBtn = Dialog.find(".DashDialogButtonSubmit") as JQuery<HTMLButtonElement>;

    CloseBtn.one("click", () => {

        HideDialog(Dialog);

    });

    RemoveDialogOnEscape(Dialog);
    SubmitDialogOnEnter(Dialog, EnterBtn);
    
}

export function HideDialog(Dialog: JQuery<HTMLDialogElement>): void {

    Dialog.focus();
    
    RelevantElements.Dialogs.Overlay.hide();
    Dialog.hide();

    Dialog.trigger("Close");
    
}

export async function WaitForDialogResponse(Dialog: JQuery<HTMLDialogElement>): Promise<boolean> {

    return new Promise((Resolve) => {

        const CloseBtn = Dialog.find(".DashDialogButtonCancel");
        const SubmitBtn = Dialog.find(".DashDialogButtonSubmit");

        CloseBtn.one("click", () => {

            Resolve(false);

        });

        SubmitBtn.one("click", () => {

            Resolve(true);
            
        });

    });

}

// Context Menu

export interface ContextMenuOption {

    Name: string,
    Icon: string,

    Action: () => void

}

export function ShowContextMenu(Event: JQuery.MouseEventBase | JQuery.TouchEventBase, Options: ContextMenuOption[]): void {
    
    // Populate Menu

    RelevantElements.CtxMenu.html("");

    Options.forEach((Option) => {

        const OptionElement = $(`
            
            <div class="DashContextMenuItem">

                ${Option.Name}
    
                ${Option.Icon}
    
            </div>
            
        `);

        OptionElement.on("click", () => {

            Option.Action();
            RelevantElements.CtxMenu.hide();

        });

        RelevantElements.CtxMenu.append(OptionElement);

    });    

    // Force render off-screen to calculate dimensions

    RelevantElements.CtxMenu.css({

        "left": "-9999px",
        "top": "-9999px",
        "display": "flex"

    });

    // Position Menu

    const MenuX = Event.pageX || 0;
    const MenuY = Event.pageY || 0;

    // Check for out of bounds

    const WindowWidth = $(window).width() || 0;
    const WindowHeight = $(window).height() || 0;

    const MenuWidth = RelevantElements.CtxMenu.width() || 0;
    const MenuHeight = RelevantElements.CtxMenu.height() || 0;

    const X = (MenuX + MenuWidth + 15) > WindowWidth ? (MenuX - MenuWidth - 15) : MenuX;
    const Y = (MenuY + MenuHeight + 15) > WindowHeight ? (MenuY - MenuHeight - 15) : MenuY;

    // Show Menu at correct position

    RelevantElements.CtxMenu.css({

        "left": X + "px",
        "top": Y + "px"

    });

    // Listen for next click (or scroll) to hide

    $(document).one("click", (E) => {

        if (!$(E.target).closest(".InlineFileActions").length) RelevantElements.CtxMenu.hide();

    });

    $(document).one("scroll", () => {

        RelevantElements.CtxMenu.hide();

    });

}
import $ from "jquery";

import { AnimationTimes, LogLevel, type Service, type SproutAccount } from "../Misc/Structs";

import { RemoveDialogOnEscape, SubmitDialogOnEnter } from "./Interactions";
import { Log, PluralizeWord } from "../Misc/Utils";

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

    Selects: {

        Services: $("#ServiceSelect"),

    },

    Containers: {

        Logs: $(".Container.ServiceLogsContainer"),

    },

    States: {

        CurrentService: {

            Label: $(".DashHeaderMessageService"),

            LogStats: {

                Warnings: $(".DashHeaderMessageCount.Warnings"),
                Errors: $(".DashHeaderMessageCount.Errors"),

            }

        }

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

// Service Select

export function PopulateServiceSelect(Services: Partial<Service>[]): void {

    // Prepend "All Services" as an option

    Services.unshift({ Name: "All Services", UID: "All" });
    
    const ServiceSelect = RelevantElements.Selects.Services;

    const Options = ServiceSelect.find("SelectMenuOptions");

    Options.html(""); // Clear existing options

    ServiceSelect.attr("Value", Services[0].UID!); // Set default value
    ServiceSelect.find(".SelectMenuLabel").text(Services[0].Name!); // Set default label

    for (const Service of Services) {

        const Option = /*html*/ `

            <SelectMenuOption Value="${Service.UID}">

                ${Service.Name}

            </SelectMenuOption>

        `;

        Options.append(Option);

    }

}

// Service State

export function SwitchToService(Service: Service): void {

    Log("Info", `Switching to service ${Service.Name} (UID ${Service.UID})`); // log the switch
        
    UpdateUIWithService(Service); // update the UI with the new service

    // Render the actual logs

    const HTML = Service.Logs.map((Log) => Log.HTML).join("");

    RelevantElements.Containers.Logs.html(HTML); // update the main panel with the new logs

    // Scroll to the bottom, which will be the most recent logs

    RelevantElements.Containers.Logs.scrollTop(RelevantElements.Containers.Logs[0].scrollHeight);

}

export function UpdateUIWithService(Service: Service): void {

    RelevantElements.States.CurrentService.Label.text(Service.Name); // update the label with the service name

    const Warnings = Service.Logs.filter((Log) => Log.Level == LogLevel.Warning).length; // get the warning count
    const Errors = Service.Logs.filter((Log) => Log.Level == LogLevel.Error).length; // get the error count

    RelevantElements.States.CurrentService.LogStats.Warnings.text(`${Warnings} ${PluralizeWord(Warnings, "Warning", "Warnings")}`); // update the warning count
    RelevantElements.States.CurrentService.LogStats.Errors.text(`${Errors} ${PluralizeWord(Errors, "Error", "Errors")}`); // update the error count

}
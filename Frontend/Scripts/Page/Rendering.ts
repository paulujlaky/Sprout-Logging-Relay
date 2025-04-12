import $ from "jquery";

import { AnimationTimes, LogLevel, type Service, type ServiceLog, type SproutAccount } from "../Misc/Structs";

import { RemoveDialogOnEscape, SubmitDialogOnEnter } from "./Interactions";
import { Log, PluralizeWord } from "../Misc/Utils";
import { GlobalStorage } from "../Main";
import { GetService } from "../Internals/Calls";

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

        Start: $(".Container.ServiceLogsStart"),
        Logs: $(".Container.ServiceLogsContent"),

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

// Footer

export function ShowFooterMessage(Mode: "Success" | "Info" | "Error", Message: string, HideTimeout: number = -1): void {

    // assign an ID lock

    const LockID = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);

    $(".Container.DashFooterMessageSection .DashFooterMessageText").html(Message);
    $(".Container.DashFooterMessageSection").css("display", "flex").attr("NotificationLockID", LockID);

    const Icons = {

        Success: `<ion-icon name="checkmark-circle-outline"></ion-icon>`,
        Info: `<ion-icon name="information-circle-outline"></ion-icon>`,
        Error: `<ion-icon name="warning-outline"></ion-icon>`

    }

    $(".Container.DashFooterMessageSection .DashFooterMessageIcon").html(Icons[Mode]);
    $(".Container.DashFooterMessageSection .DashFooterMessageIcon").css("display", "flex");
    
    if (HideTimeout != -1) {

        setTimeout(() => {

            HideFooterMessage(LockID);

        }, HideTimeout)

    };

}

export function HideFooterMessage(LockID: string): void {

    const FooterMessage = $(".Container.DashFooterMessageSection");

    if (FooterMessage.attr("NotificationLockID") == LockID) {

        FooterMessage.fadeOut(AnimationTimes.Short, () => {

            FooterMessage.css("display", "none").attr("NotificationLockID", "");

        });

    }

}

// Dialogs

export function ShowDialog(Dialog: JQuery<HTMLDialogElement>): void {

    RelevantElements.Dialogs.Overlay.show();
    Dialog.css("display", "flex");

    // Clear inputs

    Dialog.find("input").first().trigger("focus"); // only focus on the first input
    Dialog.find("input").val(""); // remove ALL values from inputs
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

                <span class="SelectMenuOptionName">
                
                    ${Service.Name}

                </span>

                    ${Service.UID !== "All" ?

                        `<span class="SelectMenuOptionDetails">

                            ${Service.UID}

                        </span>`
                        
                    : ""}

            </SelectMenuOption>

        `;

        Options.append(Option);

    }

}

// Service State

export function SwitchToService(Service: Service): void {

    Log("Info", `Switching to service ${Service.Name} (UID ${Service.UID})`); // log the switch

    GlobalStorage.CurrentService = Service; // set the current service in the global storage
        
    UpdateUIWithService(Service); // update the UI with the new service

    // Render the actual logs

    if (Service.Logs.length == 0) {

        RelevantElements.Containers.Logs.html(""); // clear the logs

        $(".ServiceLogsStart").show(); // show the start message
        return;

    } else {

        $(".ServiceLogsStart").hide(); // hide the start message

    }

    const HTML = Service.Logs.map((Log) => Log.HTML).join("");

    RelevantElements.Containers.Logs.html(HTML); // update the main panel with the new logs

    // Scroll to the bottom, which will be the most recent logs

    document.scrollingElement?.scrollTo({ // scroll to the bottom of the logs

        top: document.scrollingElement.scrollHeight,
        behavior: "smooth"

    });

}

export async function AppendLogToUI(FromServiceUID: string, Log: ServiceLog, FromSocket: boolean = false): Promise<void> {

    if (FromSocket && GlobalStorage.CurrentService?.UID != FromServiceUID && GlobalStorage.CurrentService?.UID != "All") { // remember, if current service UID is all, we want to show all logs

        return; // if the log is from a socket and the service is not the current one, ignore it

    }

    const HTML = Log.HTML; // get the HTML for the log

    if (!HTML) return; // if the HTML is empty, ignore it

    if (RelevantElements.Containers.Start.is(":visible")) { // if the start message is visible, hide it

        RelevantElements.Containers.Start.hide(); // hide the start message
        
    }

    RelevantElements.Containers.Logs.append(HTML); // append the log to the logs container

    document.scrollingElement?.scrollTo({ // scroll to the bottom of the logs

        top: document.scrollingElement.scrollHeight,
        behavior: "smooth"

    });

}

export function UpdateUIWithService(Service: Service): void {

    RelevantElements.States.CurrentService.Label.text(Service.Name); // update the label with the service name

    const Warnings = Service.Logs.filter((Log) => Log.Level == LogLevel.Warning).length; // get the warning count
    const Errors = Service.Logs.filter((Log) => Log.Level == LogLevel.Error).length; // get the error count

    RelevantElements.States.CurrentService.LogStats.Warnings.text(`${Warnings} ${PluralizeWord(Warnings, "Warning", "Warnings")}`); // update the warning count
    RelevantElements.States.CurrentService.LogStats.Errors.text(`${Errors} ${PluralizeWord(Errors, "Error", "Errors")}`); // update the error count

}
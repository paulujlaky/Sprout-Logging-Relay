import $ from 'jquery';

import { HideDialog, PopulateServiceSelect, ShowContextMenu, ShowDialog, ShowFooterMessage, SwitchToService, UpdateUIWithService, WaitForDialogResponse, type ContextMenuOption } from './Rendering';
import { CreateService, GetAllServices, GetService } from '../Internals/Calls';

const RelevantElements = {

    Document: $(document),

    Buttons: {

        AccountMenu: $(".DashHeaderAccount"),

        ReloadBtn: $("#RefreshButton"),

        CreateServiceBtn: $("#CreateService"),

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

export async function WatchPageInteractions(): Promise<void> {

    // Watchers

    HandleToggleClicks();

    HandleSelectMenuInteractions();
    HandleContextMenuInteractions();

    WatchBrandingClick();

    WatchForAccountMenu();
    WatchForReloadBtn();
    WatchForCreateService();

    PopulateServiceSelect(await GetAllServices());

    WatchForServiceSwitches(); // Must be after PopulateServiceSelect to ensure the select menu is populated

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

// Create Service Dialog

function WatchForCreateService(): void {

    RelevantElements.Buttons.CreateServiceBtn.on("click", async () => {

        const Dialog = $(".CreateServiceDialog") as JQuery<HTMLDialogElement>;

        ShowDialog(Dialog);

        const DialogResponse = await WaitForDialogResponse(Dialog);

        HideDialog(Dialog);

        if (!DialogResponse) { return; } // if the user cancelled, return

        const ServiceName = $("#ServiceName").val() as string; // get the service name from the input
        const ServiceDescription = $("#ServiceDescription").val() as string; // get the service type from the select menu
        const ServiceIconURL = $("#ServiceIconURL").val() as string; // get the service icon URL from the input

        if (!ServiceName || !ServiceDescription || !ServiceIconURL) { // if any of the values are empty, return

            ShowFooterMessage("Error", "Please complete all fields", 2_500); 

            return; // early return

        }

        const CreateResponse = await CreateService(ServiceName, ServiceDescription, ServiceIconURL); // create the service

        if (CreateResponse?.UID) {

            ShowFooterMessage("Success", `Created service ${CreateResponse.Name} (UID ${CreateResponse.UID})`, 5_000); // show success message

        } else {
            
            ShowFooterMessage("Error", "Failed to create service", 2_500); // show error message

        }

        PopulateServiceSelect(await GetAllServices()); // repopulate the service select menu

    });

}

// Reload Button

async function WatchForReloadBtn(): Promise<void> {

    RelevantElements.Buttons.ReloadBtn.on("click", async () => {

        const CurrentService = RelevantElements.Selects.ServiceSelect.attr("Value") || ""; // get the current service UID

        const Service = await GetService(CurrentService, true); // get the service (AND its logs) from the backend

        Service && SwitchToService(Service); // switch to the service
        
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

// Service Switcher/Filter

function WatchForServiceSwitches(): void {

    const SelectMenuHandler = async (Selected: SelectMenuOption) => {

        const ServiceUID = Selected.Value; // get the service UID from the select menu

        const Service = await GetService(ServiceUID, true); // get the service (AND its logs) from the backend

        if (!Service) { return; } // if the service is null, return

        SwitchToService(Service); // switch to the service
        
    }

    HandleSelectMenuResponse(RelevantElements.Selects.ServiceSelect, SelectMenuHandler); // handle the select menu response

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

export function HandleSelectMenuResponse(SelectMenu: JQuery<HTMLElement>, OptionClickCallback: (Option: SelectMenuOption) => void): void {

    const UpdateSelectMenuWithOption = (Option: SelectMenuOption): void => {

        SelectMenu.find(".SelectMenuLabel").text(Option.Name); // update the label with the new option name
        SelectMenu.attr("Value", Option.Value); // update the value of the select menu

    }

    $(document).on("click", "SelectMenuOption", (Event: JQuery.MouseEventBase | JQuery.TouchEventBase) => { // listening this way will allow for new options to be added dynamically
    
        const RelatedOption = $(Event.target).closest("SelectMenuOption");

        const OptionValue = RelatedOption.attr("Value") || ""; // get the value of the option
        const OptionName = RelatedOption.find(".SelectMenuOptionName").text(); // the element text is the name

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
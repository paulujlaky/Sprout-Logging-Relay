import type { ServiceLog, Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { AppendLogToUI, HideLoadingView, SwitchToService } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

import { GetAllServices, GetService, WatchForLogs } from "./Internals/Calls";
import { Log } from "Utils";

export const GlobalStorage: Storage = {

    User: null,

    CurrentService: null,

    Cache: {

        Services: {}

    },

    States: {
        
        Pagination: {
            
            LastOffset: 0,
            LastLimit: 25,

        }
        
    }

};

const IsAuthed = await Authorize();

if (IsAuthed) {

    const AllServices = await GetAllServices();

    WatchPageInteractions(AllServices); // Watch for page interactions

    const ConglomorateService = await GetService("All"); // Get all services with logs

    Log("Info", "All services loaded, rendering UI");

    HideLoadingView();

    ConglomorateService ? SwitchToService(ConglomorateService) : null;

    // Watch for new logs

    const LogListener = (ServiceUID: string, NewLog: ServiceLog) => {

        Log("Info", `New log from websocket for service UID ${ServiceUID}`);

        AppendLogToUI(ServiceUID, NewLog, true); // Tell it from WS so ignore if not from same service as current selected
        
    }

    WatchForLogs(LogListener);

}
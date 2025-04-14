import type { ServiceLog, Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { AppendLogToUI, HideLoadingView, SwitchToService } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

import { GetAllServices, GetService, WatchForLogs } from "./Internals/Calls";
import { Log } from "./Misc/Utils";

export const GlobalStorage: Storage = {

    User: null,

    CurrentService: null,

    Cache: {

        Services: {}

    },

    States: {
        
    }

};

const IsAuthed = await Authorize();

if (IsAuthed) {

    const AllServices = await GetAllServices();

    HideLoadingView();

    Log("Info", "All services loaded, rendering UI");

    WatchPageInteractions(AllServices); // Watch for page interactions

    const ConglomorateService = await GetService("All", true); // Get all services with logs
    
    ConglomorateService ? SwitchToService(ConglomorateService) : null;

    // Watch for new logs

    const LogListener = (ServiceUID: string, NewLog: ServiceLog) => {

        Log("Info", `New log from websocket for service UID ${ServiceUID}`);

        AppendLogToUI(ServiceUID, NewLog, true); // Tell it from WS so ignore if not from same service as current selected
        
    }

    WatchForLogs(LogListener);

}
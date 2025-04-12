import type { ServiceLog, Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { AppendLogToUI, HideLoadingView, SwitchToService } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

import { GetService, WatchForLogs } from "./Internals/Calls";

export const GlobalStorage: Storage = {

    User: null,

    CurrentService: null,

    Cache: {

        Services: {}

    }

};

const IsAuthed = await Authorize();

if (IsAuthed) {

    HideLoadingView();

    WatchPageInteractions();

    const AllServices = await GetService("All", true);
    
    AllServices ? SwitchToService(AllServices) : null;

    // Watch for new logs

    const LogListener = (ServiceUID: string, NewLog: ServiceLog) => {

        AppendLogToUI(ServiceUID, NewLog, true); // Tell it from WS so ignore if not from same service as current selected
        
    }

    WatchForLogs(LogListener);

}
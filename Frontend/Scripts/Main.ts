import type { Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { HideLoadingView, SwitchToService } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

import { GetService } from "./Internals/Calls";

export const GlobalStorage: Storage = {

    User: null,

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

}
import type { Storage } from "./Misc/Structs";

import { Authorize } from "./Systems/Auth";

import { HideLoadingView } from "./Page/Rendering";
import { WatchPageInteractions } from "./Page/Interactions";

export const GlobalStorage: Storage = {

    User: null,

    Cache: {

        Files: {},
        Dirs: {}

    }

};

const IsAuthed = await Authorize();

if (IsAuthed) {

    HideLoadingView();

    WatchPageInteractions();

}
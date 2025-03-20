import Routes from "../../Routes.json";

import $ from "jquery";

import { GlobalStorage } from "../Main";
import type { SproutAccount } from "../Misc/Structs";

import { GetURLParameter, Log, MakeRequest, RemoveURLParameter } from "../Misc/Utils";

import { UpdateAuthElements } from "../Page/Rendering";

export async function Authorize(): Promise<boolean> {

    const Token = GetURLParameter("Auth"); // In the case user has just logged in

    // Assume token is either here, or in cookies. Either way, we must POST to API/Auth to see

    const Response = await MakeRequest(Routes.Authorize, { Token: Token ? Token : null });

    if (Response?.JSON?.UID) {

        // We are authorized

        SetAccountState(Response.JSON as SproutAccount);

        return true;

    } else if (Response?.JSON?.WasFlagError) {

        $(".DashNoAccessMessage").show();

        return false;

    } else {

        window.location.href = "https://sprout.software/Accounts/Login?Redirect=SproutLoggingRelayDev"; // Redirects to sprout login page
        return false;
        
    }

}

function SetAccountState(Account: SproutAccount): void {

    Log("Info", `Welcome${CheckFirstLogin() ? "" : " back"}, ${Account.Username}!`);
        
    GlobalStorage.User = Account;
    
    RemoveURLParameter("Auth"); // Removes token from URL

    UpdateAuthElements(GlobalStorage.User);

    // Watch for log out button clicks

    $("#LogOut").on("click", async () => {

        await MakeRequest(Routes.LogOut);

        document.location.reload();

    });

}

export function CheckFirstLogin(): boolean {

    // This is written to persistent/local storage

    const WasFirst = localStorage.getItem("OOBE") == null;

    if (WasFirst) localStorage.setItem("OOBE", "true");
    
    return WasFirst;

}
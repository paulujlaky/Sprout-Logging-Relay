// File which contains most functions which utilize API calls to the backend

import type { ServiceLog, Service } from "../Misc/Structs";

import { Log, MakeRequest } from "../Misc/Utils";

import Routes from "../../Routes.json";

// Services

export async function GetAllServices(): Promise<Service[]> {

    const Response = await MakeRequest(Routes.ListServices);

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to get services list from backend.");
        return [];

    }

    return Response.JSON as Service[];

}

// Logs

export async function GetAllLogs(): Promise<ServiceLog[]> {

    const Response = await MakeRequest(Routes.ViewLogs);

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to get service logs from backend.");
        return [];

    }

    return Response.JSON as ServiceLog[]; // Will automatically be sorted by most recent

}

export async function WatchForLogs(NewLogCallback: (ServiceUID: string, Log: ServiceLog) => void): Promise<void> {

    // We will use a websocket

    const WatchLogsRouteURL = Routes.WatchForLogs.Path;
    
    const WebSocketInstance = new WebSocket(WatchLogsRouteURL);

    // Some default handlers for open/close/err events 

    WebSocketInstance.onopen = () => Log("Info", "Connected to logs websocket.");
    WebSocketInstance.onclose = () => Log("Error", "Disconnected from logs websocket.");
    WebSocketInstance.onerror = () => Log("Error", "Error in logs websocket.");

    WebSocketInstance.onmessage = (MsgEvent: MessageEvent) => {

        const Data = MsgEvent.data;

        if (Data == "Ping") {

            // Our crappy ping/pong system. For now, we handle it seperately.

            return WebSocketInstance.send("Pong");

        }

        let ParsedData = null;

        try {

            ParsedData = JSON.parse(Data);

        } catch (Err) {

            Log("Error", "Failed to parse message from websocket.");
            return;

        }

        if (!ParsedData || !ParsedData.JSON?.ServiceUID) {

            return; // May not be related data for a new log

        }

        ParsedData = ParsedData.JSON as {

            ServiceUID: string;
            Log: ServiceLog;

        }

        NewLogCallback(ParsedData.ServiceUID, ParsedData.Log); // Call the callback with the new log. here, the cb will run as many times as the logs are received

    }

}

// Misc

export async function LogOut(): Promise<void> {

    const Response = await MakeRequest(Routes.LogOut);

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to log out / release cookies.");
        return;

    }

}
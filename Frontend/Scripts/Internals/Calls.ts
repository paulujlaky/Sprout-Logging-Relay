import { GlobalStorage } from "../Main";

// File which contains most functions which utilize API calls to the backend

import type { ServiceLog, Service } from "../Misc/Structs";

import { AddItemToCache, Log, MakeRequest } from "../Misc/Utils";

import Routes from "../../Routes.json";

// Services

export async function CreateService(ServiceName: string, ServiceDescription: string, ServiceIconURL: string, NotificationThreshold: number, NotificationInterval: number): Promise<Service | null> {

    // Create a service with the given name, description, and icon URL

    const Response = await MakeRequest(Routes.CreateService, { Name: ServiceName, Description: ServiceDescription, Icon: ServiceIconURL, Threshold: NotificationThreshold, Interval: NotificationInterval });

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to create service.");
        return null;

    }

    const Service = Response.JSON as Service;

    AddItemToCache("Services", Service.UID, Service); // Add to cache

    return Service; // Will automatically be sorted by most recent

}


export async function GetAllServices(): Promise<Service[]> {

    const Response = await MakeRequest(Routes.ListServices);

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to get services list from backend.");
        return [];

    }

    const Services = Response.JSON as Service[];
    
    Services.forEach((Service) => {

        AddItemToCache("Services", Service.UID, Service); // Add to cache

    });

    return Services;

}

export async function GetService(ServiceUID: string): Promise<Service | null> {

    if (ServiceUID == "All") {

        return {

            UID: "All",

            Name: "All Services",
            Icon: "https://cdn.sprout.software/Files/Sprout/Logos/Sprout.png", 

            Description: "N/A",

            MaxLogs: 0, // Not applicable for all services

        } satisfies Service; 

    }

    // Check cache first

    let Service: Service | null = GlobalStorage.Cache.Services[ServiceUID]?.Data || null; // Check if the service is in the cache

    if (!Service) {
    
        const Services = await GetAllServices(); // If not in cache, get all services to ensure we have the latest data

        Service = Services.find((Service) => Service.UID == ServiceUID) || null; // Attempt to find the service in the list of services
        
    }

    if (!Service) {

        Log("Error", "Failed to get service from backend.");
        return null;

    }

    return Service; // Will automatically be sorted by most recent

}

// Logs

export async function GetLog(LogUID: string): Promise<ServiceLog | null> {

    const Response = await MakeRequest(Routes.GetLog, null, {}, { LogUID: LogUID });

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to get log from backend.");
        return null;

    }

    return Response.JSON as ServiceLog; 

}

export async function DeleteLog(LogUID: string): Promise<boolean> {

    const Response = await MakeRequest(Routes.DeleteLog, null, {}, { LogUID: LogUID });

    if (!Response || !Response.JSON?.Success) {

        Log("Error", "Failed to delete log from backend.");
        return false;

    }

    return true;

}

export async function GetServiceLogs(ServiceUID: string, Offset: number = 0, Limit: number = 25): Promise<ServiceLog[]> {

    // Get all logs for a service

    const Response = await MakeRequest(Routes.ViewServiceSpecificLogs, null, {}, { ServiceUID: ServiceUID }, { Offset: Offset, Limit: Limit });

    if (!Response || !Response.JSON) {

        return []; // May not be a failure... no logs

    }

    return Response.JSON as ServiceLog[]; // Will automatically be sorted by most recent

}

export async function WatchForLogs(NewLogCallback: (ServiceUID: string, Log: ServiceLog) => void): Promise<void> {

    // We will use a websocket

    const WatchLogsRouteURL = Routes.WatchForLogs.Path;
    
    let WebSocketInst = new WebSocket(WatchLogsRouteURL);

    // Some default handlers for open/close/err events

    const BindListeners = (WebSocketInstance: WebSocket) => {

        WebSocketInstance.onopen = () => Log("Info", "Connected to logs websocket.");
        WebSocketInstance.onclose = () => Log("Warning", "Disconnected from logs websocket.");
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

                // Was more likely than not just a ping/pong message, so we can ignore it safely
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

    BindListeners(WebSocketInst); // Bind the listeners to the initial instance

    setInterval(() => {

        if (WebSocketInst.readyState == WebSocket.OPEN) {

            WebSocketInst.send("Ping"); // Send a ping every 5 seconds

        }

    }, 30_000); // 30 seconds

    setInterval(() => {

        // Check if websocket is open, and if not, reconnect

        if (WebSocketInst.readyState == WebSocket.CLOSED) {

            Log("Warning", "Detected disconnection from logs websocket. Attempting to reconnect.");

            WebSocketInst = new WebSocket(WatchLogsRouteURL); // Recreate the websocket instance

            // Bind the listeners again

            BindListeners(WebSocketInst); // Rebind the listeners

        }

    }, 5_000);

}

// Misc

export async function LogOut(): Promise<void> {

    const Response = await MakeRequest(Routes.LogOut);

    if (!Response || !Response.JSON) {

        Log("Error", "Failed to log out / release cookies.");
        return;

    }

}
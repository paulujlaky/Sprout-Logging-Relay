// Type Util

export type RenderedResource<T> = { HTML: string, JSON: T };

// Internal/Frontend

export enum AnimationTimes {

    Short = 150,
    Medium = 300,
    Long = 450

}

export type CacheItem<T> = {

    Data: T;
    Expiry: number;

}

export interface GlobalCache {

    // Nothing yet

    [key: string]: { [key: string]: CacheItem<any> };

}

export interface GlobalStates {

    Scrolling: {    

        IsAtBottom: boolean;
    
    },

    Pagination: {

        LastOffset: number;
        LastLimit: number;

        IsNoMoreLogs: boolean;
        IsMakingRequests: boolean;

        Debounce: {

            Debounced: boolean;

        }

    }

}

export interface Storage {

    User: SproutAccount | null;

    CurrentService: Service | null;

    Cache: GlobalCache;

    States: GlobalStates;

}

// Auth

export interface SproutAccount {

    UUID: string;

    Username: string;
    Email: string;
    Avatar?: string;

    Flags: string[];
    
}

// Services & Logs

export enum LogLevel {

    Info = 0,
    Warning = 1,
    Error = 2

}

export type FeatureType = "code" | "link";

interface LogFeature {

    Type: FeatureType;
    Value: string;

}

interface LogContent {

    Title: string;
    Description: string;

    Features: LogFeature[];

}

export interface ServiceLog {

    ParentServiceUID: string;

    Level: LogLevel;

    Content: LogContent;

    Created: number;

    HTML: string; 

}

export interface Service {

    UID: string;

    Name: string;
    Description: string;
    Icon: string;

    MaxLogs: number;

}
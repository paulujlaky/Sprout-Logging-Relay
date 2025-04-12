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

    Level: LogLevel;

    Content: LogContent;

    Created: number;

    HTML?: string; // Is defined everywhere on backend, but was added later on... so best to be cautious

}

export interface Service {

    UID: string;

    Name: string;
    Description: string;
    Icon: string;

    Logs: ServiceLog[];

    MaxLogs: number;

}
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

export interface Storage {

    User: SproutAccount | null;
    Cache: GlobalCache;

}

// Auth

export interface SproutAccount {

    UUID: string;

    Username: string;
    Email: string;
    Avatar?: string;

    Flags: string[];
    
}
export enum UpdateAvailability {
    UpToDate = 0,
    Available,
    Unknown = -1
}

/**
 * @interface AppStatus
 * @description Represents the status of an app for the component.
 */
export interface AppStatus {
    installed: boolean
    installing: boolean
    updateAvailable: UpdateAvailability
}

export namespace model {
	
	export class ReleaseV2 {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    entityId?: number[];
	    version?: string;
	    codeVersion?: string;
	    contentVersion?: string;
	    name?: string;
	    description?: string;
	    archive: boolean;
	    app?: AppV2;
	
	    static createFrom(source: any = {}) {
	        return new ReleaseV2(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.entityId = source["entityId"];
	        this.version = source["version"];
	        this.codeVersion = source["codeVersion"];
	        this.contentVersion = source["contentVersion"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.archive = source["archive"];
	        this.app = this.convertValues(source["app"], AppV2);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ReleaseV2Batch {
	    entities?: ReleaseV2[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new ReleaseV2Batch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], ReleaseV2);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SDK {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    releases?: ReleaseV2Batch;
	
	    static createFrom(source: any = {}) {
	        return new SDK(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.releases = this.convertValues(source["releases"], ReleaseV2Batch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Presence {
	    // Go type: time
	    updatedAt?: any;
	    status?: string;
	    // Go type: uuid
	    spaceId?: any;
	    // Go type: uuid
	    serverId?: any;
	
	    static createFrom(source: any = {}) {
	        return new Presence(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.status = source["status"];
	        this.spaceId = this.convertValues(source["spaceId"], null);
	        this.serverId = this.convertValues(source["serverId"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Persona {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    name?: string;
	
	    static createFrom(source: any = {}) {
	        return new Persona(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.name = source["name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Comment {
	    id?: number[];
	    entityId?: number[];
	    userId: number[];
	    text: string;
	
	    static createFrom(source: any = {}) {
	        return new Comment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entityId = source["entityId"];
	        this.userId = source["userId"];
	        this.text = source["text"];
	    }
	}
	export class CommentBatch {
	    entities?: Comment[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new CommentBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Comment);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Likable {
	    id?: number[];
	    entityId?: number[];
	    userId: number[];
	    value: number;
	    createdAt?: Date;
	    updatedAt?: Date;
	
	    static createFrom(source: any = {}) {
	        return new Likable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entityId = source["entityId"];
	        this.userId = source["userId"];
	        this.value = source["value"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	    }
	}
	export class LikableBatch {
	    entities?: Likable[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new LikableBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Likable);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Property {
	    id?: number[];
	    entityId?: number[];
	    type: string;
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new Property(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entityId = source["entityId"];
	        this.type = source["type"];
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class PropertyBatch {
	    entities?: Property[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new PropertyBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Property);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Link {
	    id?: number[];
	    url: string;
	    name?: string;
	
	    static createFrom(source: any = {}) {
	        return new Link(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.name = source["name"];
	    }
	}
	export class LinkBatch {
	    entities?: Link[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new LinkBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Link);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class File {
	    id?: number[];
	    entityId?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    type: string;
	    url: string;
	    mime?: string;
	    size?: number;
	    version?: number;
	    deploymentType?: string;
	    platform?: string;
	    uploadedBy?: number[];
	    width?: number;
	    height?: number;
	    variation?: number;
	    originalPath?: string;
	    hash?: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entityId = source["entityId"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.type = source["type"];
	        this.url = source["url"];
	        this.mime = source["mime"];
	        this.size = source["size"];
	        this.version = source["version"];
	        this.deploymentType = source["deploymentType"];
	        this.platform = source["platform"];
	        this.uploadedBy = source["uploadedBy"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.variation = source["variation"];
	        this.originalPath = source["originalPath"];
	        this.hash = source["hash"];
	    }
	}
	export class FileBatch {
	    entities?: File[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new FileBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], File);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Accessible {
	    id?: number[];
	    entityId?: number[];
	    userId: number[];
	    username?: string;
	    isOwner: boolean;
	    canView: boolean;
	    canEdit: boolean;
	    canDelete: boolean;
	    createdAt?: Date;
	    updatedAt?: Date;
	
	    static createFrom(source: any = {}) {
	        return new Accessible(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.entityId = source["entityId"];
	        this.userId = source["userId"];
	        this.username = source["username"];
	        this.isOwner = source["isOwner"];
	        this.canView = source["canView"];
	        this.canEdit = source["canEdit"];
	        this.canDelete = source["canDelete"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	    }
	}
	export class AccessibleBatch {
	    entities?: Accessible[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new AccessibleBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], Accessible);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class User {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    email?: string;
	    apiKey?: string;
	    name?: string;
	    description?: string;
	    ip?: string;
	    geoLocation?: string;
	    isActive?: boolean;
	    isAdmin?: boolean;
	    isMuted?: boolean;
	    isBanned?: boolean;
	    isInternal?: boolean;
	    // Go type: time
	    lastSeenAt?: any;
	    // Go type: time
	    activatedAt?: any;
	    allowEmails?: boolean;
	    experience?: number;
	    level?: number;
	    rank?: string;
	    ethAddress?: string;
	    address?: string;
	    defaultPersonaId?: number[];
	    // Go type: Persona
	    defaultPersona?: any;
	    // Go type: Presence
	    presence?: any;
	    isEmailConfirmed?: boolean;
	    isAddressConfirmed?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.email = source["email"];
	        this.apiKey = source["apiKey"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.ip = source["ip"];
	        this.geoLocation = source["geoLocation"];
	        this.isActive = source["isActive"];
	        this.isAdmin = source["isAdmin"];
	        this.isMuted = source["isMuted"];
	        this.isBanned = source["isBanned"];
	        this.isInternal = source["isInternal"];
	        this.lastSeenAt = this.convertValues(source["lastSeenAt"], null);
	        this.activatedAt = this.convertValues(source["activatedAt"], null);
	        this.allowEmails = source["allowEmails"];
	        this.experience = source["experience"];
	        this.level = source["level"];
	        this.rank = source["rank"];
	        this.ethAddress = source["ethAddress"];
	        this.address = source["address"];
	        this.defaultPersonaId = source["defaultPersonaId"];
	        this.defaultPersona = this.convertValues(source["defaultPersona"], null);
	        this.presence = this.convertValues(source["presence"], null);
	        this.isEmailConfirmed = source["isEmailConfirmed"];
	        this.isAddressConfirmed = source["isAddressConfirmed"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppV2 {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    name?: string;
	    description?: string;
	    external: boolean;
	    sdk?: SDK;
	    releases?: ReleaseV2Batch;
	
	    static createFrom(source: any = {}) {
	        return new AppV2(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.external = source["external"];
	        this.sdk = this.convertValues(source["sdk"], SDK);
	        this.releases = this.convertValues(source["releases"], ReleaseV2Batch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppV2Batch {
	    entities?: AppV2[];
	    offset?: number;
	    limit?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new AppV2Batch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entities = this.convertValues(source["entities"], AppV2);
	        this.offset = source["offset"];
	        this.limit = source["limit"];
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class LauncherV2 {
	    id?: number[];
	    createdAt?: Date;
	    updatedAt?: Date;
	    entityType?: string;
	    public?: boolean;
	    views?: number;
	    // Go type: User
	    owner?: any;
	    // Go type: AccessibleBatch
	    accessibles?: any;
	    // Go type: FileBatch
	    files?: any;
	    // Go type: LinkBatch
	    links?: any;
	    // Go type: PropertyBatch
	    properties?: any;
	    // Go type: LikableBatch
	    likables?: any;
	    // Go type: CommentBatch
	    comments?: any;
	    liked?: number;
	    likes?: number;
	    dislikes?: number;
	    name: string;
	    description: string;
	    releases?: ReleaseV2Batch;
	    apps?: AppV2Batch;
	
	    static createFrom(source: any = {}) {
	        return new LauncherV2(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.createdAt = new Date(source["createdAt"]);
	        this.updatedAt = new Date(source["updatedAt"]);
	        this.entityType = source["entityType"];
	        this.public = source["public"];
	        this.views = source["views"];
	        this.owner = this.convertValues(source["owner"], null);
	        this.accessibles = this.convertValues(source["accessibles"], null);
	        this.files = this.convertValues(source["files"], null);
	        this.links = this.convertValues(source["links"], null);
	        this.properties = this.convertValues(source["properties"], null);
	        this.likables = this.convertValues(source["likables"], null);
	        this.comments = this.convertValues(source["comments"], null);
	        this.liked = source["liked"];
	        this.likes = source["likes"];
	        this.dislikes = source["dislikes"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.releases = this.convertValues(source["releases"], ReleaseV2Batch);
	        this.apps = this.convertValues(source["apps"], AppV2Batch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}


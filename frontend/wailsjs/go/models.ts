export namespace gui {
	
	export class BackupInfo {
	    id: string;
	    timestamp: string;
	    items: string[];
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new BackupInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = source["timestamp"];
	        this.items = source["items"];
	        this.size = source["size"];
	    }
	}
	export class BackupPreview {
	    items: string[];
	
	    static createFrom(source: any = {}) {
	        return new BackupPreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = source["items"];
	    }
	}
	export class CleanerInfo {
	    name: string;
	    items: string[];
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.items = source["items"];
	        this.count = source["count"];
	    }
	}
	export class WipeMethodInfo {
	    id: number;
	    name: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new WipeMethodInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	    }
	}
	export class WiperInfo {
	    totalSpace: number;
	    freeSpace: number;
	    volume: string;
	    methods: WipeMethodInfo[];
	
	    static createFrom(source: any = {}) {
	        return new WiperInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalSpace = source["totalSpace"];
	        this.freeSpace = source["freeSpace"];
	        this.volume = source["volume"];
	        this.methods = this.convertValues(source["methods"], WipeMethodInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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


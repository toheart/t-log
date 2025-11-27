export namespace config {
	
	export class AppConfig {
	    root_path: string;
	    hotkey: string;
	    history_days: number;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.root_path = source["root_path"];
	        this.hotkey = source["hotkey"];
	        this.history_days = source["history_days"];
	    }
	}

}

export namespace note {
	
	export class NoteEntry {
	    content: string;
	    timestamp: string;
	    date: string;
	
	    static createFrom(source: any = {}) {
	        return new NoteEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.timestamp = source["timestamp"];
	        this.date = source["date"];
	    }
	}

}


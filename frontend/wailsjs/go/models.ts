export namespace command {
	
	export class Command {
	    id: string;
	    title: string;
	    description: string;
	    usage: string;
	
	    static createFrom(source: any = {}) {
	        return new Command(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.usage = source["usage"];
	    }
	}

}

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
	
	export class DailyNote {
	    date: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new DailyNote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.content = source["content"];
	    }
	}
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
	export class SearchResult {
	    content: string;
	    date: string;
	    time: string;
	    filePath: string;
	    lineNo: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.date = source["date"];
	        this.time = source["time"];
	        this.filePath = source["filePath"];
	        this.lineNo = source["lineNo"];
	    }
	}

}


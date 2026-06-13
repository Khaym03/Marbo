export namespace domain {
	
	export class Response {
	    text: string;
	    zone_id?: string;
	    suggestions?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.text = source["text"];
	        this.zone_id = source["zone_id"];
	        this.suggestions = source["suggestions"];
	    }
	}

}

export namespace runtime {
	
	export class ClarificationOption {
	    intent_id: string;
	    label: string;
	
	    static createFrom(source: any = {}) {
	        return new ClarificationOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.intent_id = source["intent_id"];
	        this.label = source["label"];
	    }
	}
	export class ClarificationData {
	    options: ClarificationOption[];
	
	    static createFrom(source: any = {}) {
	        return new ClarificationData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.options = this.convertValues(source["options"], ClarificationOption);
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
	
	export class ConfidenceData {
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new ConfidenceData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.score = source["score"];
	    }
	}
	export class Runtime {
	
	
	    static createFrom(source: any = {}) {
	        return new Runtime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class TraceData {
	
	
	    static createFrom(source: any = {}) {
	        return new TraceData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class RuntimeExtension {
	    confidence?: ConfidenceData;
	    clarify?: ClarificationData;
	    // Go type: TraceData
	    trace?: any;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeExtension(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.confidence = this.convertValues(source["confidence"], ConfidenceData);
	        this.clarify = this.convertValues(source["clarify"], ClarificationData);
	        this.trace = this.convertValues(source["trace"], null);
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
	export class RuntimeResult {
	    type: string;
	    response: domain.Response;
	    intent_id: string;
	    zone_id: string;
	    flow_id: string;
	    node_id: string;
	    extension?: RuntimeExtension;
	
	    static createFrom(source: any = {}) {
	        return new RuntimeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.response = this.convertValues(source["response"], domain.Response);
	        this.intent_id = source["intent_id"];
	        this.zone_id = source["zone_id"];
	        this.flow_id = source["flow_id"];
	        this.node_id = source["node_id"];
	        this.extension = this.convertValues(source["extension"], RuntimeExtension);
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


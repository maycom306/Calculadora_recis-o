export namespace rescisao {
	
	export class Servidor {
	    nome: string;
	    matricula: string;
	    cargo: string;
	
	    static createFrom(source: any = {}) {
	        return new Servidor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nome = source["nome"];
	        this.matricula = source["matricula"];
	        this.cargo = source["cargo"];
	    }
	}
	export class DadosRescisao {
	    servidor: Servidor;
	    banco: string;
	    agencia: string;
	    conta: string;
	    pix: string;
	    dataAdmissao: string;
	    dataDesligamento: string;
	    diasTrabalhados: number;
	    salarioBaseCentavos: number;
	    avosDecimo: number;
	    avosDecimoDisponiveis: number;
	    avosDecimoDevidos: number;
	    avosFeriasPagos: number;
	    avosFeriasDisponiveis: number;
	    avosFeriasDevidos: number;
	
	    static createFrom(source: any = {}) {
	        return new DadosRescisao(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.servidor = this.convertValues(source["servidor"], Servidor);
	        this.banco = source["banco"];
	        this.agencia = source["agencia"];
	        this.conta = source["conta"];
	        this.pix = source["pix"];
	        this.dataAdmissao = source["dataAdmissao"];
	        this.dataDesligamento = source["dataDesligamento"];
	        this.diasTrabalhados = source["diasTrabalhados"];
	        this.salarioBaseCentavos = source["salarioBaseCentavos"];
	        this.avosDecimo = source["avosDecimo"];
	        this.avosDecimoDisponiveis = source["avosDecimoDisponiveis"];
	        this.avosDecimoDevidos = source["avosDecimoDevidos"];
	        this.avosFeriasPagos = source["avosFeriasPagos"];
	        this.avosFeriasDisponiveis = source["avosFeriasDisponiveis"];
	        this.avosFeriasDevidos = source["avosFeriasDevidos"];
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
	export class EntradaRescisao {
	    nome: string;
	    matricula: string;
	    cargo: string;
	    banco: string;
	    agencia: string;
	    conta: string;
	    pix: string;
	    dataAdmissao: string;
	    dataDesligamento: string;
	    salarioBase: string;
	    avosDecimo: string;
	    avosFeriasPagos: string;
	
	    static createFrom(source: any = {}) {
	        return new EntradaRescisao(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.nome = source["nome"];
	        this.matricula = source["matricula"];
	        this.cargo = source["cargo"];
	        this.banco = source["banco"];
	        this.agencia = source["agencia"];
	        this.conta = source["conta"];
	        this.pix = source["pix"];
	        this.dataAdmissao = source["dataAdmissao"];
	        this.dataDesligamento = source["dataDesligamento"];
	        this.salarioBase = source["salarioBase"];
	        this.avosDecimo = source["avosDecimo"];
	        this.avosFeriasPagos = source["avosFeriasPagos"];
	    }
	}
	export class ResultadoRescisao {
	    decimoTerceiroCentavos: number;
	    feriasCentavos: number;
	    tercoFeriasCentavos: number;
	    diasTrabalhadosCentavos: number;
	    avosFeriasDevidos: number;
	    totalCentavos: number;
	
	    static createFrom(source: any = {}) {
	        return new ResultadoRescisao(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.decimoTerceiroCentavos = source["decimoTerceiroCentavos"];
	        this.feriasCentavos = source["feriasCentavos"];
	        this.tercoFeriasCentavos = source["tercoFeriasCentavos"];
	        this.diasTrabalhadosCentavos = source["diasTrabalhadosCentavos"];
	        this.avosFeriasDevidos = source["avosFeriasDevidos"];
	        this.totalCentavos = source["totalCentavos"];
	    }
	}

}


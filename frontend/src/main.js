import './style.css';
import './app.css';

import {CalcularResultado, GerarPDF, ValidarDados} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
        <header class="topbar">
            <div class="brand-mark" aria-hidden="true">R</div>
            <div class="brand-copy">
                <span class="eyebrow">RECURSOS HUMANOS</span>
                <h1>RescisaoDesk</h1>
            </div>
            <span class="local-status"><span></span> Aplicativo local</span>
        </header>
        <main class="workspace">
            <section class="form-panel" aria-labelledby="page-title">
                <div class="section-heading">
                    <div>
                        <p class="eyebrow">NOVO DOCUMENTO</p>
                        <h2 id="page-title">Dados do servidor</h2>
                    </div>
                    <span class="step-indicator"><span id="current-step">01</span> <span>/ 04</span></span>
                </div>
                <form id="term-form" novalidate>
                        <section id="server-step" class="wizard-step" aria-labelledby="page-title">
                            <div class="field-grid">
                                <label class="field field-wide">
                                    <span>Nome completo <b>*</b></span>
                                    <input name="nome" type="text" autocomplete="name" maxlength="160" required />
                                </label>
                                <label class="field">
                                    <span>Matrícula</span>
                                    <input name="matricula" type="text" autocomplete="off" maxlength="40" />
                                </label>
                                <label class="field">
                                    <span>Cargo</span>
                                    <input name="cargo" type="text" autocomplete="organization-title" maxlength="120" />
                                </label>
                        </div>
                            <div class="form-actions">
                                <p class="required-note"><b>*</b> Campos obrigatórios</p>
                                <button id="next-step" class="primary-button" type="button">
                                    <span>Avançar</span><span aria-hidden="true">→</span>
                                </button>
                            </div>
                        </section>

                        <section id="termination-step" class="wizard-step" aria-labelledby="termination-title" hidden>
                            <div class="field-grid">
                                <label class="field">
                                    <span>Data de admissão <b>*</b></span>
                                    <input name="dataAdmissao" type="date" required />
                                </label>
                                <label class="field">
                                    <span>Data de desligamento <b>*</b></span>
                                    <input name="dataDesligamento" type="date" required />
                                </label>
                                <label class="field">
                                    <span>Salário / base (R$) <b>*</b></span>
                                    <input name="salarioBase" type="text" inputmode="decimal" placeholder="0,00" autocomplete="off" required />
                                </label>
                                <label class="field">
                                    <span>Avos do 13º já recebidos <b>*</b></span>
                                    <input name="avosDecimo" type="number" min="0" step="1" value="0" required />
                                </label>
                                <label class="field">
                                    <span>Avos de férias já gozados/pagos <b>*</b></span>
                                    <input name="avosFeriasPagos" type="number" min="0" step="1" value="0" required />
                                </label>
                            </div>
                            <div class="form-actions">
                                <button id="previous-step" class="secondary-button" type="button">Voltar</button>
                                <button id="go-to-bank" class="primary-button" type="button">
                                    <span>Próximo</span><span aria-hidden="true">→</span>
                                </button>
                            </div>
                        </section>

                        <section id="bank-step" class="wizard-step" aria-labelledby="bank-title" hidden>
                            <div class="field-grid">
                                <label class="field">
                                    <span>Banco</span>
                                    <input name="banco" type="text" maxlength="80" placeholder="Ex.: Banco do Brasil" />
                                </label>
                                <label class="field">
                                    <span>Agência</span>
                                    <input name="agencia" type="text" maxlength="20" placeholder="Ex.: 1234-5" />
                                </label>
                                <label class="field">
                                    <span>Conta bancária</span>
                                    <input name="conta" type="text" maxlength="30" placeholder="Ex.: 123456-7" />
                                </label>
                                <label class="field field-wide">
                                    <span>Chave Pix</span>
                                    <input name="pix" type="text" maxlength="120" placeholder="CPF, e-mail, telefone ou chave aleatória" />
                                </label>
                            </div>
                            <div class="form-actions">
                                <button id="back-to-termination" class="secondary-button" type="button">Voltar</button>
                                <button id="validate-data" class="primary-button" type="submit">
                                    <span>Ver cálculo</span><span aria-hidden="true">→</span>
                                </button>
                            </div>
                        </section>
                        <section id="calculation-step" class="wizard-step" aria-labelledby="calculation-title" hidden>
                            <div class="calculation-heading">
                                <p class="eyebrow">VERBAS INCLUÍDAS</p>
                                <h2 id="calculation-title">Conferência do cálculo</h2>
                                <p class="calculation-caption">O 13º considera os avos adquiridos desde a admissão, descontados os já recebidos.</p>
                            </div>
                            <div class="calculation-list" role="list">
                                <div class="calculation-row informational-row" role="listitem">
                                    <span><strong>Dias no mês do desligamento</strong><small>Quantidade de dias considerada</small></span>
                                    <output id="days-worked-count">--</output>
                                </div>
                                <div class="calculation-row deduction-row" role="listitem">
                                    <span><strong>Valor dos dias trabalhados</strong><small>Exibido separadamente; não compõe o total</small></span>
                                    <output id="days-worked-deduction">Pendente</output>
                                </div>
                                <div class="calculation-row" role="listitem">
                                    <span><strong>13º proporcional</strong><small id="thirteenth-avos-summary">Avos aguardando validação</small></span>
                                    <output id="thirteenth-value">Pendente</output>
                                </div>
                                <div class="calculation-row" role="listitem">
                                    <span><strong>Férias proporcionais</strong><small id="vacation-avos-summary">Avos aguardando validação</small></span>
                                    <output id="vacation-value">Pendente</output>
                                </div>
                                <div class="calculation-row" role="listitem">
                                    <span><strong>1/3 de férias</strong><small>Sobre as férias proporcionais</small></span>
                                    <output id="vacation-third-value">Pendente</output>
                                </div>
                            </div>
                            <div class="calculation-total">
                                <strong>Total das verbas (sem dias trabalhados)</strong>
                                <output id="calculation-total-value">Pendente</output>
                            </div>
                            <p class="calculation-blocker" role="note">Os dias trabalhados são apresentados à parte e não entram no total das verbas.</p>
                            <div class="form-actions">
                                <button id="back-to-bank" class="secondary-button" type="button">Voltar</button>
                                <button id="generate-pdf" class="primary-button" type="button">
                                    <span>Gerar PDF</span><span aria-hidden="true">→</span>
                                </button>
                            </div>
                        </section>
                    <p id="form-status" class="form-status" role="status" aria-live="polite"></p>
                </form>
            </section>
            <aside id="validation-panel" class="validation-panel" hidden aria-live="polite">
                <div class="validation-icon" aria-hidden="true">✓</div>
                <p class="eyebrow">DADOS CONFERIDOS</p>
                <h2>Entrada validada</h2>
                <dl id="validated-data" class="validated-data"></dl>
            </aside>
        </main>
        <footer class="app-footer">
            <span>RescisaoDesk</span>
            <span>Os dados desta sessão não são armazenados.</span>
        </footer>
`;

const form = document.querySelector('#term-form');
const statusElement = document.querySelector('#form-status');
const validationPanel = document.querySelector('#validation-panel');
const validatedData = document.querySelector('#validated-data');
const serverStep = document.querySelector('#server-step');
const terminationStep = document.querySelector('#termination-step');
const bankStep = document.querySelector('#bank-step');
const calculationStep = document.querySelector('#calculation-step');
const currentStep = document.querySelector('#current-step');
const pageTitle = document.querySelector('#page-title');
const nameInput = form.elements.namedItem('nome');

document.querySelector('#next-step').addEventListener('click', () => {
    statusElement.textContent = '';
    statusElement.classList.remove('is-error');
    validationPanel.hidden = true;

    if (!nameInput.reportValidity()) return;

    serverStep.hidden = true;
    terminationStep.hidden = false;
    calculationStep.hidden = true;
    currentStep.textContent = '02';
    pageTitle.textContent = 'Dados da rescisão';
    terminationStep.querySelector('input').focus();
});

document.querySelector('#previous-step').addEventListener('click', () => {
    statusElement.textContent = '';
    statusElement.classList.remove('is-error');
    validationPanel.hidden = true;
    terminationStep.hidden = true;
    calculationStep.hidden = true;
    serverStep.hidden = false;
    currentStep.textContent = '01';
    pageTitle.textContent = 'Dados do servidor';
    nameInput.focus();
});

document.querySelector('#go-to-bank').addEventListener('click', () => {
    terminationStep.hidden = true;
    bankStep.hidden = false;
    calculationStep.hidden = true;
    currentStep.textContent = '03';
    pageTitle.textContent = 'Dados bancários';
    bankStep.querySelector('input').focus();
});

document.querySelector('#back-to-termination').addEventListener('click', () => {
    bankStep.hidden = true;
    terminationStep.hidden = false;
    currentStep.textContent = '02';
    pageTitle.textContent = 'Dados da rescisão';
    terminationStep.querySelector('input').focus();
});

document.querySelector('#back-to-bank').addEventListener('click', () => {
    calculationStep.hidden = true;
    bankStep.hidden = false;
    currentStep.textContent = '03';
    pageTitle.textContent = 'Dados bancários';
    bankStep.querySelector('input').focus();
});

const formatMoney = (centavos) => new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
}).format(centavos / 100);

form.addEventListener('submit', async (event) => {
    event.preventDefault();
    if (bankStep.hidden) return;

    statusElement.textContent = '';
    statusElement.classList.remove('is-error');
    validationPanel.hidden = true;

    const formData = Object.fromEntries(new FormData(form).entries());
    const submitButton = document.querySelector('#validate-data');
    submitButton.disabled = true;

    try {
        const result = await ValidarDados(formData);
        const calculation = await CalcularResultado(formData);
        document.querySelector('#days-worked-count').textContent = `${result.diasTrabalhados} ${result.diasTrabalhados === 1 ? 'dia' : 'dias'}`;
        document.querySelector('#days-worked-deduction').textContent = formatMoney(calculation.diasTrabalhadosCentavos);
        document.querySelector('#thirteenth-value').textContent = formatMoney(calculation.decimoTerceiroCentavos);
        document.querySelector('#vacation-value').textContent = formatMoney(calculation.feriasCentavos);
        document.querySelector('#vacation-third-value').textContent = formatMoney(calculation.tercoFeriasCentavos);
        document.querySelector('#calculation-total-value').textContent = formatMoney(calculation.totalCentavos);
        const formatDate = (value) => {
            const [year, month, day] = value.split('-');
            return `${day}/${month}/${year}`;
        };
        const rows = [
            ['Servidor', result.servidor.nome],
            ['Matrícula', result.servidor.matricula || 'Não informada'],
            ['Cargo', result.servidor.cargo || 'Não informado'],
            ['Banco', result.banco || 'Não informado'],
            ['Agência', result.agencia || 'Não informada'],
            ['Conta bancária', result.conta || 'Não informada'],
            ['Pix', result.pix || 'Não informado'],
            ['Admissão', formatDate(result.dataAdmissao)],
            ['Desligamento', formatDate(result.dataDesligamento)],
            ['Salário / base', formData.salarioBase],
            ['Avos de 13º adquiridos', String(result.avosDecimoDisponiveis)],
            ['Avos do 13º já recebidos', String(result.avosDecimo)],
            ['Avos de 13º ainda devidos', String(result.avosDecimoDevidos)],
            ['Avos de férias disponíveis', String(result.avosFeriasDisponiveis)],
            ['Avos de férias já gozados/pagos', String(result.avosFeriasPagos)],
            ['Avos de férias ainda devidos', String(result.avosFeriasDevidos)],
        ];

        validatedData.replaceChildren(...rows.flatMap(([label, value]) => {
            const term = document.createElement('dt');
            const description = document.createElement('dd');
            term.textContent = label;
            description.textContent = value;
            return [term, description];
        }));
        document.querySelector('#vacation-avos-summary').textContent =
            `${result.avosFeriasDisponiveis} disponíveis - ${result.avosFeriasPagos} pagos = ${result.avosFeriasDevidos} devidos`;
        document.querySelector('#thirteenth-avos-summary').textContent =
            `${result.avosDecimoDisponiveis} adquiridos - ${result.avosDecimo} recebidos = ${result.avosDecimoDevidos} devidos`;
        validationPanel.hidden = false;
        bankStep.hidden = true;
        terminationStep.hidden = true;
        calculationStep.hidden = false;
        currentStep.textContent = '04';
        pageTitle.textContent = 'Cálculo da rescisão';
        statusElement.textContent = 'Cálculo concluído. Confira os valores e gere o PDF.';
    } catch (error) {
        statusElement.textContent = error instanceof Error ? error.message : String(error);
        statusElement.classList.add('is-error');
    } finally {
        submitButton.disabled = false;
    }
});

document.querySelector('#generate-pdf').addEventListener('click', async (event) => {
    const button = event.currentTarget;
    const formData = Object.fromEntries(new FormData(form).entries());
    button.disabled = true;
    statusElement.textContent = 'Gerando e abrindo o PDF...';
    statusElement.classList.remove('is-error');

    try {
        const path = await GerarPDF(formData);
        statusElement.textContent = `PDF gerado e aberto: ${path}`;
    } catch (error) {
        statusElement.textContent = error instanceof Error ? error.message : String(error);
        statusElement.classList.add('is-error');
    } finally {
        button.disabled = false;
    }
});

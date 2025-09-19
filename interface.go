// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	

	"github.com/nsf/termbox-go"
)

// Define um tipo Cor para encapsuladar as cores do termbox
type Cor = termbox.Attribute


// Definições de cores utilizadas no jogo
const (
	CorPadrao     Cor = termbox.ColorDefault
	CorCinzaEscuro    = termbox.ColorDarkGray
	CorVermelho       = termbox.ColorRed
	CorVerde          = termbox.ColorGreen
	CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede    = termbox.ColorDarkGray
	CorTexto          = termbox.ColorDarkGray
	CorAmarelo        = termbox.ColorYellow //Adicionado (cor do patrulheiro)
)

// EventoTeclado representa uma ação detectada do teclado (como mover, sair ou interagir)
type EventoTeclado struct {
	Tipo  string // "sair", "interagir", "mover"
	Tecla rune   // Tecla pressionada, usada no caso de movimento
}

// Modificação: Comando para renderização thread-safe
type ComandoRender struct {
	Dados    *Jogo
	Resposta chan bool
}
//Modificação: Canal de comandos de renderização
var renderChan = make(chan ComandoRender, 10)

// Inicializa a interface gráfica usando termbox
func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	//Modificação: Inicia o controlador de interface, evita condição de corrida
	go interfaceControlador()
}

// Encerra o uso da interface termbox
func interfaceFinalizar() {
	termbox.Close()
}

// Modificação: Controlador de renderização - processa todas as operações de desenho sequencialmente
func interfaceControlador() {
	for cmd := range renderChan {
		interfaceLimparTela()

		// Desenha todos os elementos do mapa
		for y, linha := range cmd.Dados.Mapa {
			for x, elem := range linha {
				interfaceDesenharElemento(x, y, elem)
			}
		}

		// Desenha o personagem sobre o mapa
		interfaceDesenharElemento(cmd.Dados.PosX, cmd.Dados.PosY, Personagem)

		// Desenha a barra de status
		interfaceDesenharBarraDeStatus(cmd.Dados)

		// Força a atualização do terminal
		interfaceAtualizarTela()
		
		// Confirma que o desenho foi concluído
		cmd.Resposta <- true
	}
}

// Lê um evento do teclado e o traduz para um EventoTeclado
func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

// Modificação: Renderiza todo o estado atual do jogo na tela de forma thread-safe
func interfaceDesenharJogo(jogo *Jogo) {
	resposta := make(chan bool)
	renderChan <- ComandoRender{
		Dados:    jogo,
		Resposta: resposta,
	}
	<-resposta // Aguarda conclusão do desenho
}

// Limpa a tela do terminal
func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

// Força a atualização da tela do terminal com os dados desenhados
func interfaceAtualizarTela() {
	termbox.Flush()
}

// Desenha um elemento na posição (x, y)
func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

// Exibe uma barra de status com informações úteis ao jogador
func interfaceDesenharBarraDeStatus(jogo *Jogo) {
	// Linha de status dinâmica
	for i, c := range jogo.StatusMsg {
		termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
	}

	// Instruções fixas
	msg := "Use WASD para mover e E para interagir. ESC para sair."
	for i, c := range msg {
		termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
	}
}


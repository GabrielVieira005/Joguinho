package main

import (
    "math/rand"
    "time"
)

type MovimentoPatrulheiro struct {
    DX, DY int
}

type Posicao struct {
    X, Y int
}

type PatrulheiroState struct {
    posX, posY int
    under      Elemento
    movChan    chan MovimentoPatrulheiro
    ackChan    chan Posicao
}

var canalMorte = make(chan bool)

// Inicia um patrulheiro com canais próprios
func iniciarPatrulheiro(x, y int, jogo *Jogo) {
    state := &PatrulheiroState{
        posX:    x,
        posY:    y,
        under:   Vazio, // começa sobre chão vazio
        movChan: make(chan MovimentoPatrulheiro),
        ackChan: make(chan Posicao),
    }

    

    // Goroutine de movimento aleatória
    go func() {
        direcoes := []struct{ dx, dy int }{
            {1, 0}, {-1, 0}, {0, 1}, {0, -1},
        }
        for {
            dir := direcoes[rand.Intn(len(direcoes))]
            state.movChan <- MovimentoPatrulheiro{DX: dir.dx, DY: dir.dy}
            <-state.ackChan
            time.Sleep(400 * time.Millisecond)
        }
    }()

    // Goroutina processadora (dona das alterações no mapa para este patruleiro)
    go func() {
        for {
            mov := <-state.movChan
            nx, ny := state.posX+mov.DX, state.posY+mov.DY

            if jogoPodeMoverPara(jogo, nx, ny) {
                // Restaura o que estava antes na posição atual
                jogo.Mapa[state.posY][state.posX] = state.under

                // Guarda o que existe no destino
                nextUnder := jogo.Mapa[ny][nx]

                // Coloca o patrulheiro no destino
                jogo.Mapa[ny][nx] = Patrulheiro

                // Atualiza estado
                state.posX, state.posY = nx, ny
                state.under = nextUnder

                interfaceDesenharJogo(jogo)
            }

            // Confirma posição
            state.ackChan <- Posicao{X: state.posX, Y: state.posY}

            // Verifica proximidade do jogador
            if dist2(state.posX, state.posY, jogo.PosX, jogo.PosY) <= 4 {
                canalMorte <- true
            }
        }
    }()
}

// Distância ao quadrado (evita sqrt)
func dist2(x1, y1, x2, y2 int) int {
    dx := x2 - x1
    dy := y2 - y1
    return dx*dx + dy*dy
}

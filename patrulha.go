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

var (
    canalPatrulhaMov   = make(chan MovimentoPatrulheiro)
    canalPatrulhaAck   = make(chan Posicao)
    canalPatrulhaInit  = make(chan Posicao)
    canalMorte         = make(chan bool)
)

// Goroutine que gera movimento aleatório contínuo
func iniciarPatrulheiro(x, y int) {
    go func() {
        canalPatrulhaInit <- Posicao{X: x, Y: y}

        direcoes := []struct{ dx, dy int }{
            {1, 0}, {-1, 0}, {0, 1}, {0, -1},
        }

        for {
            dir := direcoes[rand.Intn(len(direcoes))]
            canalPatrulhaMov <- MovimentoPatrulheiro{DX: dir.dx, DY: dir.dy}
            <-canalPatrulhaAck // só para sincronizar
            time.Sleep(400 * time.Millisecond)
        }
    }()
}

// Processador que aplica movimentos e verifica proximidade
func executarProcessadorPatrulha(jogo *Jogo) {
    var cur Posicao
    under := Vazio

    cur = <-canalPatrulhaInit

    for {
        mov := <-canalPatrulhaMov
        nx, ny := cur.X+mov.DX, cur.Y+mov.DY

        if jogoPodeMoverPara(jogo, nx, ny) {
            jogo.Mapa[cur.Y][cur.X] = under
            nextUnder := jogo.Mapa[ny][nx]
            jogo.Mapa[ny][nx] = Patrulheiro
            cur = Posicao{X: nx, Y: ny}
            under = nextUnder
            interfaceDesenharJogo(jogo)
        }

        canalPatrulhaAck <- cur

        if dist2(cur.X, cur.Y, jogo.PosX, jogo.PosY) <= 4 {
            canalMorte <- true
        }
    }
}

// Distância ao quadrado (evita sqrt)
func dist2(x1, y1, x2, y2 int) int {
    dx := x2 - x1
    dy := y2 - y1
    return dx*dx + dy*dy
}

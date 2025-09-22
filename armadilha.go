package main

import (
	"time"
)

var armadilhaAtiva = make(map[[2]int]bool)
var canalArmadilha = make(chan bool)

// Inicia uma armadilha temporizada com canal e timeout
func iniciarArmadilha(x, y int, jogo *Jogo) {
	pos := [2]int{x, y}
	if armadilhaAtiva[pos] {
		return // já existe armadilha ativa nesta posição
	}
	armadilhaAtiva[pos] = true

	go func() {
		defer func() {
			// Recupera de qualquer panic para evitar crash
			if r := recover(); r != nil {
				armadilhaAtiva[pos] = false
			}
		}()
		
		for {
			// Ativa a armadilha
			jogo.Mapa[y][x] = Armadilha
			interfaceDesenharJogo(jogo)

			// Canal para detectar ativação
			ativou := make(chan struct{})
			done := make(chan struct{})

			// Goroutine para saber se o jogador tá na armadilha
			go func() {
				defer func() {
					if r := recover(); r != nil {// Se houver panic, apenas fecha o canal
					}
					close(done)
				}()
				
				for {
					time.Sleep(100 * time.Millisecond)
					
					if jogo.Mapa[y][x].simbolo != Armadilha.simbolo {
						return
					}
					
					if jogo.PosX == x && jogo.PosY == y {
						select {
						case ativou <- struct{}{}:
						default:
						}
						return
					}
				}
			}()

			select {
			case <-ativou:
				canalArmadilha <- true
				jogo.PresoAte = time.Now().Add(3 * time.Second) // prende por 3 segundos
				jogo.Mapa[y][x] = ArmadilhaUsada
				interfaceDesenharJogo(jogo)
				time.Sleep(1 * time.Second)
			case <-time.After(5 * time.Second):
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				time.Sleep(1 * time.Second)
			}
			<-done
			
			// Aguarda 4 segundos antes de reativar a armadilha
			time.Sleep(4 * time.Second)
			
			// Só reativa se o local ainda estiver vazio ou com armadilha usada
			elemento := jogo.Mapa[y][x]
			if elemento.simbolo != Vazio.simbolo && elemento.simbolo != ArmadilhaUsada.simbolo {
				break
			}
		}
		armadilhaAtiva[pos] = false
	}()
}
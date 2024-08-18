package game

import (
	"fmt"
	"math/rand"
	"push-pon/display"
	"push-pon/keyboard"
	"time"
)

type GameController struct {
	level         int
	successCount  int
	totalScore    int
	remainingKeys map[int]bool
	keyboard      *keyboard.KeyboardController
	display       *display.DisplayController
	startTime     time.Time
	gameOver      bool
	gameStarted   bool
	ledArray      []bool // 光らせるスイッチの状態を保持
}

func NewGameController(kb *keyboard.KeyboardController, display *display.DisplayController) *GameController {
	return &GameController{
		level:       1,
		keyboard:    kb,
		display:     display,
		gameOver:    false,
		gameStarted: false,
	}
}

func (game *GameController) StartGame() {
	game.gameStarted = true
	game.display.ShowMessage("Starting Game!")
	time.Sleep(1 * time.Second)
	game.StartNextTurn()
}

func (game *GameController) IsGameStarted() bool {
	return game.gameStarted
}

func (game *GameController) StartNextTurn() {
	if game.gameOver || !game.gameStarted {
		return
	}

	// すべてのLEDを消灯してリセット
	game.keyboard.DisableAllKeys()

	// LED配列の初期化と光らせるキーの決定
	game.remainingKeys = make(map[int]bool)
	numKeys := game.level
	totalKeys := len(game.keyboard.RowPins()) * len(game.keyboard.ColPins()) // マトリックスキーボードの全キー数

	game.ledArray = make([]bool, totalKeys)

	// 光らせるキーの選択
	for i := 0; i < numKeys; i++ {
		keyIndex := rand.Intn(totalKeys)
		for game.ledArray[keyIndex] {
			keyIndex = rand.Intn(totalKeys) // 同じキーを選択しないように再選択
		}
		game.ledArray[keyIndex] = true
		game.remainingKeys[keyIndex] = true
	}

	// LEDの点灯（緑色）
	game.keyboard.LightKeys(game.ledArray, Green)

	// ディスプレイに現在のレベルを表示
	game.display.ShowMessage(fmt.Sprintf("Level %d", game.level))
	game.startTime = time.Now()
}

func (game *GameController) CheckKeyPress() bool {
	if game.gameOver || !game.gameStarted {
		return false
	}

	keyIndex, pressed := game.keyboard.ScanKeys()
	if pressed {
		// インデックスが範囲内であることを確認し、正しいキーが押されたか判定
		if keyIndex >= 0 && keyIndex < len(game.ledArray) && game.remainingKeys[keyIndex] {
			// 正しいキーが押された場合
			game.remainingKeys[keyIndex] = false
			game.keyboard.TurnOffKey(keyIndex)
			delete(game.remainingKeys, keyIndex)

			// すべてのキーが押されたら次のターンへ
			if len(game.remainingKeys) == 0 {
				elapsed := time.Since(game.startTime).Milliseconds()
				game.totalScore += int(elapsed)
				game.successCount++

				// レベルアップ処理
				if game.successCount >= 5 && game.level == 10 {
					game.display.ShowResult(game.totalScore, true)
					game.EndGame()
					return true
				}

				if game.successCount >= 5 {
					game.level++
					game.successCount = 0
				}

				// 次のターンへ
				game.StartNextTurn()
				return true
			}
		} else {
			// 間違ったキーが押された場合は赤色で表示
			game.keyboard.LightKeys([]bool{true}, Red)
			game.display.ShowMessage("Wrong key! Game Over!")
			game.display.ShowResult(game.totalScore, false)
			game.EndGame()
			return true
		}
	}

	return false
}

func (game *GameController) EndGame() {
	game.gameOver = true
	game.keyboard.DisableAllKeys() // ゲーム終了時にすべてのLEDを消灯
}

func (game *GameController) ResetGame() {
	game.level = 1
	game.successCount = 0
	game.totalScore = 0
	game.gameOver = false
	game.gameStarted = false
	game.StartNextTurn()
}

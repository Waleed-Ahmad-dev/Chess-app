package server

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chess Game</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }

        .container {
            display: flex;
            gap: 30px;
            max-width: 1400px;
            width: 100%;
            align-items: flex-start;
        }

        .game-area {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 30px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        }

        .board-container {
            position: relative;
        }

        .board {
            display: grid;
            grid-template-columns: repeat(8, 70px);
            grid-template-rows: repeat(8, 70px);
            border: 4px solid #333;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        }

        .square {
            width: 70px;
            height: 70px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 48px;
            cursor: pointer;
            transition: all 0.2s;
            position: relative;
        }

        .square.light {
            background-color: #f0d9b5;
        }

        .square.dark {
            background-color: #b58863;
        }

        .square.selected {
            background-color: #7fc97f !important;
            box-shadow: inset 0 0 0 3px #2d7a2d;
        }

        .square.legal-move {
            position: relative;
        }

        .square.legal-move::after {
            content: '';
            position: absolute;
            width: 20px;
            height: 20px;
            background-color: rgba(0, 128, 0, 0.4);
            border-radius: 50%;
            pointer-events: none;
        }

        .square.legal-move.has-piece::after {
            width: 100%;
            height: 100%;
            border-radius: 0;
            background-color: rgba(255, 0, 0, 0.3);
        }

        .square.last-move {
            background-color: rgba(255, 255, 0, 0.3) !important;
        }

        .square:hover {
            filter: brightness(1.1);
        }

        .sidebar {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 20px;
            padding: 30px;
            min-width: 300px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        }

        .status-panel {
            margin-bottom: 30px;
        }

        .status-panel h2 {
            color: #333;
            margin-bottom: 15px;
            font-size: 24px;
        }

        .turn-indicator {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 15px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 10px;
            color: white;
            font-size: 20px;
            font-weight: bold;
            margin-bottom: 15px;
        }

        .turn-dot {
            width: 20px;
            height: 20px;
            border-radius: 50%;
            border: 2px solid white;
        }

        .turn-dot.white {
            background-color: white;
        }

        .turn-dot.black {
            background-color: #333;
        }

        .alert {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 15px;
            font-weight: bold;
        }

        .alert.check {
            background-color: #ff6b6b;
            color: white;
        }

        .alert.checkmate {
            background-color: #51cf66;
            color: white;
        }

        .alert.stalemate {
            background-color: #ffd43b;
            color: #333;
        }

        .captured-pieces {
            margin-bottom: 30px;
        }

        .captured-pieces h3 {
            color: #333;
            margin-bottom: 10px;
        }

        .captured-list {
            display: flex;
            flex-wrap: wrap;
            gap: 5px;
            min-height: 40px;
            padding: 10px;
            background: #f8f9fa;
            border-radius: 8px;
        }

        .captured-piece {
            font-size: 24px;
            opacity: 0.7;
        }

        .controls {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }

        button {
            padding: 15px 25px;
            font-size: 16px;
            font-weight: bold;
            border: none;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.3s;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #6c757d;
            color: white;
        }

        .btn-secondary:hover {
            background: #5a6268;
            transform: translateY(-2px);
        }

        .move-history {
            margin-top: 20px;
        }

        .move-history h3 {
            color: #333;
            margin-bottom: 10px;
        }

        .history-list {
            max-height: 200px;
            overflow-y: auto;
            background: #f8f9fa;
            border-radius: 8px;
            padding: 10px;
        }

        .move-item {
            padding: 5px 10px;
            margin: 2px 0;
            background: white;
            border-radius: 4px;
            font-family: monospace;
        }

        .promotion-dialog {
            display: none;
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
            z-index: 1000;
        }

        .promotion-dialog.active {
            display: block;
        }

        .promotion-dialog h3 {
            margin-bottom: 20px;
            color: #333;
        }

        .promotion-pieces {
            display: flex;
            gap: 15px;
        }

        .promotion-piece {
            width: 80px;
            height: 80px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 60px;
            background: #f0d9b5;
            border-radius: 10px;
            cursor: pointer;
            transition: all 0.2s;
        }

        .promotion-piece:hover {
            background: #b58863;
            transform: scale(1.1);
        }

        .overlay {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.5);
            z-index: 999;
        }

        .overlay.active {
            display: block;
        }

        @media (max-width: 1200px) {
            .container {
                flex-direction: column;
                align-items: center;
            }
        }

        @media (max-width: 600px) {
            .board {
                grid-template-columns: repeat(8, 45px);
                grid-template-rows: repeat(8, 45px);
            }
            
            .square {
                width: 45px;
                height: 45px;
                font-size: 32px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="game-area">
            <div class="board-container">
                <div class="board" id="board"></div>
            </div>
        </div>

        <div class="sidebar">
            <div class="status-panel">
                <h2>Game Status</h2>
                <div class="turn-indicator" id="turnIndicator">
                    <div class="turn-dot white" id="turnDot"></div>
                    <span id="turnText">White's Turn</span>
                </div>
                <div id="alerts"></div>
            </div>

            <div class="captured-pieces">
                <h3>Captured Pieces</h3>
                <div class="captured-list" id="capturedList"></div>
            </div>

            <div class="controls">
                <button class="btn-primary" onclick="newGame()">New Game</button>
                <button class="btn-secondary" onclick="undoMove()">Undo Move</button>
            </div>

            <div class="move-history">
                <h3>Move History</h3>
                <div class="history-list" id="historyList"></div>
            </div>
        </div>
    </div>

    <div class="overlay" id="overlay"></div>
    <div class="promotion-dialog" id="promotionDialog">
        <h3>Choose Promotion Piece</h3>
        <div class="promotion-pieces" id="promotionPieces"></div>
    </div>

    <script>
        let sessionId = null;
        let gameState = null;
        let selectedSquare = null;
        let pendingPromotion = null;

        const pieceUnicode = {
            'P': '♙', 'N': '♘', 'B': '♗', 'R': '♖', 'Q': '♕', 'K': '♔',
            'p': '♟', 'n': '♞', 'b': '♝', 'r': '♜', 'q': '♛', 'k': '♚',
            '.': ''
        };

        async function newGame() {
            try {
                const response = await fetch('/api/new-game', { method: 'POST' });
                const data = await response.json();
                sessionId = data.sessionId;
                gameState = data.state;
                renderGame();
            } catch (error) {
                console.error('Error starting new game:', error);
            }
        }

        async function makeMove(from, to, promotion = '') {
            try {
                const move = indexToCoord(from) + indexToCoord(to) + promotion;
                const response = await fetch('/api/make-move', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ sessionId, move })
                });
                const data = await response.json();
                
                if (data.success) {
                    gameState = data.state;
                    selectedSquare = null;
                    renderGame();
                } else {
                    alert('Invalid move: ' + data.error);
                }
            } catch (error) {
                console.error('Error making move:', error);
            }
        }

        function renderGame() {
            if (!gameState) return;

            renderBoard();
            renderStatus();
            renderCapturedPieces();
        }

        function renderBoard() {
            const board = document.getElementById('board');
            board.innerHTML = '';

            for (let rank = 7; rank >= 0; rank--) {
                for (let file = 0; file < 8; file++) {
                    const index = rank * 8 + file;
                    const square = document.createElement('div');
                    square.className = 'square';
                    square.className += (rank + file) % 2 === 0 ? ' light' : ' dark';
                    square.dataset.index = index;

                    const piece = gameState.board[index];
                    square.textContent = pieceUnicode[piece] || '';

                    if (piece !== '.') {
                        square.classList.add('has-piece');
                    }

                    if (selectedSquare === index) {
                        square.classList.add('selected');
                    }

                    if (gameState.lastMove) {
                        const from = coordToIndex(gameState.lastMove.substring(0, 2));
                        const to = coordToIndex(gameState.lastMove.substring(2, 4));
                        if (index === from || index === to) {
                            square.classList.add('last-move');
                        }
                    }

                    if (selectedSquare !== null) {
                        const moveStr = indexToCoord(selectedSquare) + indexToCoord(index);
                        if (gameState.legalMoves.some(m => m.startsWith(moveStr))) {
                            square.classList.add('legal-move');
                        }
                    }

                    square.onclick = () => handleSquareClick(index);
                    board.appendChild(square);
                }
            }
        }

        function handleSquareClick(index) {
            if (gameState.gameOver) return;

            const piece = gameState.board[index];
            const isOwnPiece = (gameState.turn === 'White' && piece === piece.toUpperCase() && piece !== '.') ||
                               (gameState.turn === 'Black' && piece === piece.toLowerCase() && piece !== '.');

            if (selectedSquare === null) {
                if (isOwnPiece) {
                    selectedSquare = index;
                    renderBoard();
                }
            } else {
                const moveStr = indexToCoord(selectedSquare) + indexToCoord(index);
                const possibleMoves = gameState.legalMoves.filter(m => m.startsWith(moveStr));

                if (possibleMoves.length > 0) {
                    if (possibleMoves.length === 1) {
                        makeMove(selectedSquare, index);
                    } else {
                        // Promotion
                        pendingPromotion = { from: selectedSquare, to: index };
                        showPromotionDialog();
                    }
                } else if (isOwnPiece) {
                    selectedSquare = index;
                    renderBoard();
                } else {
                    selectedSquare = null;
                    renderBoard();
                }
            }
        }

        function showPromotionDialog() {
            const dialog = document.getElementById('promotionDialog');
            const overlay = document.getElementById('overlay');
            const pieces = document.getElementById('promotionPieces');
            
            pieces.innerHTML = '';
            const promotionPieces = gameState.turn === 'White' ? 
                ['Q', 'R', 'B', 'N'] : ['q', 'r', 'b', 'n'];
            const promotionLetters = ['q', 'r', 'b', 'n'];

            promotionPieces.forEach((piece, i) => {
                const div = document.createElement('div');
                div.className = 'promotion-piece';
                div.textContent = pieceUnicode[piece];
                div.onclick = () => {
                    makeMove(pendingPromotion.from, pendingPromotion.to, promotionLetters[i]);
                    hidePromotionDialog();
                };
                pieces.appendChild(div);
            });

            dialog.classList.add('active');
            overlay.classList.add('active');
        }

        function hidePromotionDialog() {
            document.getElementById('promotionDialog').classList.remove('active');
            document.getElementById('overlay').classList.remove('active');
            pendingPromotion = null;
        }

        function renderStatus() {
            const turnText = document.getElementById('turnText');
            const turnDot = document.getElementById('turnDot');
            const alerts = document.getElementById('alerts');

            turnText.textContent = gameState.turn + "'s Turn";
            turnDot.className = 'turn-dot ' + gameState.turn.toLowerCase();

            alerts.innerHTML = '';

            if (gameState.gameOver) {
                const alert = document.createElement('div');
                if (gameState.isStalemate) {
                    alert.className = 'alert stalemate';
                    alert.textContent = '🤝 Stalemate! Draw!';
                } else {
                    alert.className = 'alert checkmate';
                    alert.textContent = '🎉 Checkmate! ' + gameState.winner + ' wins!';
                }
                alerts.appendChild(alert);
            } else if (gameState.inCheck) {
                const alert = document.createElement('div');
                alert.className = 'alert check';
                alert.textContent = '⚠️ Check!';
                alerts.appendChild(alert);
            }
        }

        function renderCapturedPieces() {
            const list = document.getElementById('capturedList');
            list.innerHTML = '';
            
            gameState.capturedPieces.forEach(piece => {
                const span = document.createElement('span');
                span.className = 'captured-piece';
                span.textContent = pieceUnicode[piece.substring(1)];
                list.appendChild(span);
            });
        }

        function indexToCoord(index) {
            const file = String.fromCharCode(97 + (index % 8));
            const rank = Math.floor(index / 8) + 1;
            return file + rank;
        }

        function coordToIndex(coord) {
            const file = coord.charCodeAt(0) - 97;
            const rank = parseInt(coord[1]) - 1;
            return rank * 8 + file;
        }

        function undoMove() {
            alert('Undo feature coming soon!');
        }

        // Start a new game on load
        newGame();
    </script>
</body>
</html>
`

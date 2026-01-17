package server

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chess - Monochrome Glass</title>
    <style>
        :root {
            --bg-dark: #000000;
            --bg-light: #1a1a1a;
            --glass-surface: rgba(255, 255, 255, 0.03);
            --glass-border: rgba(255, 255, 255, 0.1);
            --glass-highlight: rgba(255, 255, 255, 0.08);
            --text-main: #ffffff;
            --text-muted: #888888;
            
            --sq-light: rgba(255, 255, 255, 0.1);
            --sq-dark: rgba(0, 0, 0, 0.3);
            --sq-hover: rgba(255, 255, 255, 0.2);
            --sq-selected: rgba(255, 255, 255, 0.25);
            --sq-last-move: rgba(255, 255, 255, 0.15);
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', 'Segoe UI', sans-serif;
            background: radial-gradient(circle at 50% 50%, #2a2a2a 0%, #000000 120%);
            min-height: 100vh;
            color: var(--text-main);
            display: flex;
            justify-content: center;
            align-items: center;
            overflow: hidden; /* Prevent scroll on full screen apps */
        }

        .container {
            display: flex;
            gap: 40px;
            max-width: 1400px;
            width: 95%;
            height: 90vh;
            align-items: center;
            justify-content: center;
        }

        /* --- Glassmorphism Panels --- */
        .glass-panel {
            background: var(--glass-surface);
            backdrop-filter: blur(20px);
            -webkit-backdrop-filter: blur(20px);
            border: 1px solid var(--glass-border);
            border-radius: 24px;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
        }

        .game-area {
            padding: 40px;
            display: flex;
            justify-content: center;
            align-items: center;
        }

        .sidebar {
            padding: 30px;
            min-width: 350px;
            max-width: 400px;
            height: 100%;
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        /* --- Chess Board --- */
        .board {
            display: grid;
            grid-template-columns: repeat(8, 75px);
            grid-template-rows: repeat(8, 75px);
            border: 1px solid var(--glass-border);
            border-radius: 4px;
            overflow: hidden;
            box-shadow: 0 10px 40px rgba(0,0,0,0.6);
        }

        .square {
            width: 75px;
            height: 75px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 52px;
            cursor: pointer;
            position: relative;
            transition: background-color 0.2s ease;
            user-select: none;
        }

        .square.light { background-color: var(--sq-light); }
        .square.dark { background-color: var(--sq-dark); }

        .square:hover {
            background-color: var(--sq-hover);
        }

        .square.selected {
            background-color: var(--sq-selected) !important;
            box-shadow: inset 0 0 15px rgba(255,255,255,0.1);
        }

        .square.last-move {
            background-color: var(--sq-last-move) !important;
        }

        /* --- Pieces Styling --- */
        /* White pieces: Pure white with a subtle shadow */
        .piece-white {
            color: #ffffff;
            text-shadow: 0 4px 10px rgba(0,0,0,0.5);
            filter: drop-shadow(0 0 2px rgba(0,0,0,0.8));
        }

        /* Black pieces: Dark grey with a white outline/glow to be visible on black */
        .piece-black {
            color: #1a1a1a;
            text-shadow: 0 0 2px rgba(255,255,255,0.4); 
        }

        /* --- Legal Move Indicators --- */
        .square.legal-move::after {
            content: '';
            position: absolute;
            width: 14px;
            height: 14px;
            background-color: rgba(255, 255, 255, 0.2);
            border-radius: 50%;
            pointer-events: none;
            box-shadow: 0 0 5px rgba(0,0,0,0.5);
        }

        .square.legal-move:hover::after {
            background-color: rgba(255, 255, 255, 0.5);
        }

        .square.legal-move.has-piece::after {
            width: 100%;
            height: 100%;
            border-radius: 0;
            background: none;
            box-shadow: inset 0 0 0 4px rgba(255, 255, 255, 0.2);
        }

        /* --- Typography & UI Elements --- */
        h2, h3 {
            font-weight: 300;
            letter-spacing: 1px;
            text-transform: uppercase;
            font-size: 14px;
            color: var(--text-muted);
            margin-bottom: 10px;
            border-bottom: 1px solid var(--glass-border);
            padding-bottom: 5px;
        }

        .status-header {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 20px;
            letter-spacing: -0.5px;
            color: var(--text-main);
        }

        .turn-indicator {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 15px;
            background: rgba(255,255,255,0.02);
            border-radius: 12px;
            border: 1px solid var(--glass-border);
        }

        .turn-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
        }
        .turn-dot.white { background: #fff; box-shadow: 0 0 10px rgba(255,255,255,0.5); }
        .turn-dot.black { background: #000; border: 1px solid #555; }

        /* --- Alerts (Check/Mate) --- */
        .alert {
            margin-top: 15px;
            padding: 12px;
            border-radius: 8px;
            font-weight: 600;
            text-align: center;
            font-size: 14px;
            letter-spacing: 1px;
            border: 1px solid currentColor;
            background: rgba(0,0,0,0.3);
        }
        .alert.check { color: #ffffff; border-color: #ffffff; }
        .alert.checkmate { color: #ffffff; background: #ffffff; color: #000; }
        .alert.stalemate { color: #aaaaaa; border-color: #aaaaaa; }

        /* --- Captured Pieces --- */
        .captured-list {
            display: flex;
            flex-wrap: wrap;
            gap: 5px;
            min-height: 30px;
            margin-bottom: 20px;
        }
        .captured-piece {
            font-size: 20px;
            opacity: 0.6;
            color: var(--text-main);
        }

        /* --- Buttons --- */
        .controls {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 10px;
            margin-top: auto;
        }

        button {
            padding: 15px;
            background: transparent;
            border: 1px solid var(--glass-border);
            color: var(--text-main);
            font-family: inherit;
            font-weight: 500;
            text-transform: uppercase;
            font-size: 12px;
            letter-spacing: 1px;
            cursor: pointer;
            border-radius: 8px;
            transition: all 0.2s;
        }

        button:hover {
            background: rgba(255, 255, 255, 0.1);
            border-color: rgba(255, 255, 255, 0.3);
        }

        button.btn-primary {
            background: rgba(255, 255, 255, 0.9);
            color: black;
            border: none;
        }
        button.btn-primary:hover {
            background: #ffffff;
            box-shadow: 0 0 15px rgba(255,255,255,0.3);
        }

        /* --- History --- */
        .history-list {
            flex: 1;
            overflow-y: auto;
            background: rgba(0, 0, 0, 0.2);
            border-radius: 12px;
            padding: 10px;
            border: 1px solid var(--glass-border);
            font-family: 'Roboto Mono', monospace;
            font-size: 13px;
        }
        
        .move-row {
            display: flex;
            justify-content: space-between;
            padding: 4px 8px;
            border-bottom: 1px solid rgba(255,255,255,0.05);
        }
        .move-row:last-child { border-bottom: none; }
        .move-num { color: var(--text-muted); width: 30px; }
        .move-white { color: var(--text-main); }
        .move-black { color: var(--text-muted); }

        /* --- Scrollbar --- */
        ::-webkit-scrollbar { width: 6px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.2); border-radius: 3px; }

        /* --- Promotion Dialog --- */
        .overlay {
            position: fixed;
            top: 0; left: 0; width: 100%; height: 100%;
            background: rgba(0, 0, 0, 0.8);
            backdrop-filter: blur(5px);
            z-index: 999;
            display: none;
            opacity: 0;
            transition: opacity 0.3s;
        }
        .overlay.active { display: block; opacity: 1; }

        .promotion-dialog {
            position: fixed;
            top: 50%; left: 50%;
            transform: translate(-50%, -50%) scale(0.9);
            background: #111;
            border: 1px solid #333;
            padding: 30px;
            border-radius: 20px;
            z-index: 1000;
            display: none;
            transition: transform 0.3s;
            box-shadow: 0 20px 60px rgba(0,0,0,0.8);
            text-align: center;
        }
        .promotion-dialog.active { display: block; transform: translate(-50%, -50%) scale(1); }

        .promotion-pieces {
            display: flex;
            gap: 20px;
            margin-top: 20px;
        }
        .promotion-piece {
            width: 70px; height: 70px;
            display: flex; align-items: center; justify-content: center;
            font-size: 40px;
            background: #222;
            border-radius: 12px;
            cursor: pointer;
            border: 1px solid #333;
            color: #fff;
            transition: all 0.2s;
        }
        .promotion-piece:hover { background: #333; border-color: #fff; }

        @media (max-width: 1000px) {
            .container { flex-direction: column; height: auto; padding: 20px; }
            .sidebar { width: 100%; max-width: none; }
            .board { grid-template-columns: repeat(8, 1fr); grid-template-rows: repeat(8, 1fr); width: 100vw; max-width: 90vw; height: 90vw; max-height: 90vw; }
            .square { width: auto; height: auto; font-size: 8vw; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="glass-panel game-area">
            <div class="board" id="board"></div>
        </div>

        <div class="glass-panel sidebar">
            <div>
                <div class="status-header">Chess</div>
                <div class="turn-indicator">
                    <div class="turn-dot white" id="turnDot"></div>
                    <span id="turnText">White's Turn</span>
                </div>
                <div id="alerts"></div>
            </div>

            <div>
                <h3>Captured</h3>
                <div class="captured-list" id="capturedList"></div>
            </div>

            <div class="history-list" id="historyList">
                </div>

            <div class="controls">
                <button class="btn-secondary" onclick="undoMove()">Undo</button>
                <button class="btn-primary" onclick="newGame()">New Game</button>
            </div>
        </div>
    </div>

    <div class="overlay" id="overlay"></div>
    <div class="promotion-dialog" id="promotionDialog">
        <h3>Promote Pawn</h3>
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
            } catch (error) { console.error(error); }
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
            } catch (error) { console.error(error); }
        }

        async function undoMove() {
             try {
                const response = await fetch('/api/undo-move', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ sessionId })
                });
                const data = await response.json();
                if (data.success) {
                    gameState = data.state;
                    selectedSquare = null;
                    renderGame();
                }
            } catch (error) { console.error(error); }
        }

        function renderGame() {
            if (!gameState) return;
            renderBoard();
            renderStatus();
            renderCapturedPieces();
            renderHistory();
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
                    
                    const piece = gameState.board[index];
                    square.textContent = pieceUnicode[piece] || '';

                    // Add Specific Class for White/Black pieces to control contrast
                    if (piece !== '.') {
                        square.classList.add('has-piece');
                        if (piece === piece.toUpperCase()) {
                            square.classList.add('piece-white');
                        } else {
                            square.classList.add('piece-black');
                        }
                    }

                    // Selection & Highlights
                    if (selectedSquare === index) square.classList.add('selected');

                    if (gameState.lastMove) {
                        const from = coordToIndex(gameState.lastMove.substring(0, 2));
                        const to = coordToIndex(gameState.lastMove.substring(2, 4));
                        if (index === from || index === to) square.classList.add('last-move');
                    }

                    // Legal Moves
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
                    selectedSquare = index; // Switch selection
                    renderBoard();
                } else {
                    selectedSquare = null; // Deselect
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
                    alert.textContent = 'STALEMATE';
                } else {
                    alert.className = 'alert checkmate';
                    alert.textContent = 'CHECKMATE - ' + gameState.winner.toUpperCase() + ' WINS';
                }
                alerts.appendChild(alert);
            } else if (gameState.inCheck) {
                const alert = document.createElement('div');
                alert.className = 'alert check';
                alert.textContent = 'CHECK';
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

        function renderHistory() {
            const list = document.getElementById('historyList');
            // This is a simple history list. 
            // In a real app, you'd parse PGN or store moves more structured in Go.
            // For now, we don't have the full move history list passed from backend cleanly
            // other than the 'LastMove'. 
            // If you want full history, you'd need to update the Go struct to return []string History.
            // For now, I'll leave this empty to keep the UI clean, or you can update the backend.
            list.innerHTML = '<div style="color:#555; text-align:center; padding-top:20px;">History logs...</div>';
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

        newGame();
    </script>
</body>
</html>
`

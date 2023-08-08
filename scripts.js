const messageContainer = document.getElementById("message-container");
const nicknameInput = document.getElementById("nickname");
const guessInput = document.getElementById("guess");

// WebSocket bağlantısı için değişkenler
let id;
let roomId;
const serverAddr = "ws://localhost:8181/ws";
const ws = new WebSocket(serverAddr);

// WebSocket olayları dinleme
ws.onmessage = function (event) {
    const message = event.data;
    const response = JSON.parse(message);

    if (response.cmd === "join") {
        // join işleminin başarılı olup olmadığını kontrol etmek için kullanabilirsiniz
    } else if (response.event === "joinedRoom") {
        roomId = response.room;
        handleJoinedRoom();
    } else if (response.event === "gameOver") {
        handleGameOver(response.rankings);
    }
};

ws.onerror = function (event) {
    showMessage("Sunucuyla bağlantı hatası!");
};

// Yeni bir oyun başlatma fonksiyonları
function joinGame() {
    const nickname = nicknameInput.value;
    if (nickname === "") {
        document.getElementById("nickname-error").style.display = "block";
        return;
    }

    document.getElementById("nickname-error").style.display = "none";

    showLoading("loading");

    fetch(`http://localhost:1234/register?nickname=${encodeURIComponent(nickname)}`)
        .then((response) => response.json())
        .then((data) => {
            id = data.id;
            const request = {
                cmd: "join",
                id: id,
                nickname: nickname,
            };
            ws.send(JSON.stringify(request));
        })
        .catch((error) => {});
}

// Oyuna katılma işlemi için ayrı bir fonksiyon ekledik
function joinGame2() {
    const nickname = document.getElementById("nickname").value;
    if (nickname === "") {
        document.getElementById("nickname-error").style.display = "block";
        return;
    }

    showLoading("new-game-loading");

    const request = {
        cmd: "join",
        id: id,
    };
    ws.send(JSON.stringify(request));
}

// Tahmin yapma fonksiyonu
function guess() {
    const guessValue = guessInput.value;

    const request = {
        cmd: "guess",
        id: id,
        room: roomId,
        data: parseInt(guessValue),
    };

    ws.send(JSON.stringify(request));
}

// Oyuna katılma işlemi başarılı olduğunda çalışacak fonksiyon
function handleJoinedRoom() {
    const container = document.getElementById("container");
    container.style.display = "none";

    const gameOverContainer = document.getElementById("game-over-container");
    gameOverContainer.style.display = "none";

    const welcomeContainer = document.getElementById("welcome-container");
    welcomeContainer.style.display = "";
}

// Oyun sonuçları için ekranda mesaj gösteren fonksiyon
function handleGameOver(rankings) {
    const welcomeContainer = document.getElementById("welcome-container");
    welcomeContainer.style.display = "none";

    const gameOverContainer = document.getElementById("game-over-container");
    gameOverContainer.style.display = "";

    for (let i = 0; i < 3; i++) {
        if (id === rankings[i].player) {
            const winnerMessage = rankings[i].rank === 1 ? "Kazandınız!" : "Kaybettiniz!";
            showGameOverMessage(winnerMessage, rankings[i].rank);
            break;
        }
    }

    hideLoading();
}

// Mesajları ekrana yazdıran fonksiyon
function showMessage(message) {
    const messageElement = document.createElement("div");
    messageElement.textContent = message;
    messageContainer.appendChild(messageElement);
    messageContainer.scrollTop = messageContainer.scrollHeight;
}

// Yükleniyor animasyonunu gösteren fonksiyon
function showLoading(loadingName) {
    const loadingDiv = document.getElementById(loadingName);
    const loader = loadingDiv.querySelector(".loader");

    loadingDiv.style.display = "flex";
    loader.style.animation = "spin 1s linear infinite";
}

// Oyun sonuçlarını gösteren fonksiyon
function showGameOverMessage(winner, winningOrder) {
    const kazanmaMesajiElement = document.getElementById("kazanma-mesaji");
    const kazanmaSirasiElement = document.getElementById("kazanma-sirasi");
    const trophyImg = document.querySelector(".trophy");
    trophyImg.style.display = "";

    if (winner === "Kazandınız!") {
        trophyImg.src = "trophy.jpg"; // Kazanan için farklı bir kupa resmi kullanabilirsiniz
    } else {
        trophyImg.style.display = "none";
    }

    kazanmaMesajiElement.textContent = winner;
    kazanmaSirasiElement.textContent = "Kazanma Sırası: " + winningOrder + ".";
}

// Yükleniyor animasyonunu gizleyen fonksiyon
function hideLoading() {
    const loadingDiv = document.getElementById("new-game-loading");
    const loader = loadingDiv.querySelector(".loader");

    loader.style.animation = "";
    loadingDiv.style.display = "none";
}

//TODO gizli sayıyı oyun bitince göster
//TODO server a bağlantıda problem olunca bağlantı problemi oldu sayfayı yenile yazdır
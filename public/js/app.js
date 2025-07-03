const emojis = [
  "😀",
  "😂",
  "😍",
  "🤔",
  "😎",
  "😭",
  "😡",
  "🤯",
  "😴",
  "🥳",
  "🥺",
  "🤠",
  "🤡",
  "👽",
  "👾",
  "🎃",
  "👻",
  "🤖",
  "🧠",
  "🦷",
  "👀",
  "👂",
  "👃",
  "👅",
  "👄",
  "💋",
  "💘",
  "💝",
  "💖",
  "💗",
  "💓",
  "💞",
  "💕",
  "💔",
  "❤️",
  "🧡",
  "💛",
  "💚",
  "💙",
  "💜",
  "🤎",
  "🖤",
  "🤍",
  "��",
  "💢",
  "💥",
  "💫",
  "💦",
  "💨",
  "🕳️",
  "💣",
  "💬",
  "👁️‍🗨️",
  "🗨️",
  "🗯️",
  "💭",
  "🚀",
];

const setupScreen = document.getElementById("setup-screen");
const gameScreen = document.getElementById("game-screen");
const gameBoard = document.getElementById("game-board");
const movesSpan = document.getElementById("moves");
const winModal = document.getElementById("win-modal");
const finalMovesSpan = document.getElementById("final-moves");
const menuBtn = document.getElementById("menu-btn");
const sizeButtons = document.querySelectorAll(".size-buttons button");

const leaderboardBtn = document.getElementById("leaderboard-btn");
const leaderboardScreen = document.getElementById("leaderboard-screen");
const leaderboardTable = document
  .getElementById("leaderboard-table")
  .getElementsByTagName("tbody")[0];
const backToMenuBtn = document.getElementById("back-to-menu-btn");
const leaderboardForm = document.getElementById("leaderboard-form");
const initialsInput = document.getElementById("initials-input");
const leaderboardFilterButtons = document.querySelectorAll(
  ".leaderboard-filters button",
);

let moves = 0;
let flippedCards = [];
let matchedPairs = 0;
let totalPairs = 0;
let currentBoardSize = 0;

const LEADERBOARD_KEY = "emojiMatchLeaderboard";

function getLeaderboard() {
  // Returns the entire leaderboard object, e.g., { "2": [...], "4": [...] }
  return JSON.parse(localStorage.getItem(LEADERBOARD_KEY)) || {};
}

function saveLeaderboard(leaderboard) {
  localStorage.setItem(LEADERBOARD_KEY, JSON.stringify(leaderboard));
}

function addScore(size, initials, score) {
  const leaderboard = getLeaderboard();
  if (!leaderboard[size]) {
    leaderboard[size] = [];
  }
  const sizeLeaderboard = leaderboard[size];
  sizeLeaderboard.push({ initials, score });
  sizeLeaderboard.sort((a, b) => a.score - b.score);
  leaderboard[size] = sizeLeaderboard.slice(0, 10); // Keep top 10 for the specific size
  saveLeaderboard(leaderboard);
}

function displayLeaderboard(size) {
  // Update active button state
  leaderboardFilterButtons.forEach((button) => {
    if (parseInt(button.dataset.size, 10) === size) {
      button.classList.add("active");
    } else {
      button.classList.remove("active");
    }
  });

  const allScores = getLeaderboard();
  const sizeScores = allScores[size] || [];

  leaderboardTable.innerHTML = "";
  if (sizeScores.length === 0) {
    const row = leaderboardTable.insertRow();
    const cell = row.insertCell();
    cell.colSpan = 3;
    cell.textContent = "No scores yet for this size!";
  } else {
    sizeScores.forEach((entry, index) => {
      const row = leaderboardTable.insertRow();
      row.insertCell().textContent = index + 1;
      row.insertCell().textContent = entry.initials;
      row.insertCell().textContent = entry.score;
    });
  }
}

function shuffle(array) {
  let currentIndex = array.length,
    randomIndex;
  while (currentIndex != 0) {
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex--;
    [array[currentIndex], array[randomIndex]] = [
      array[randomIndex],
      array[currentIndex],
    ];
  }
  return array;
}

function getBoardDimensions(size) {
  switch (size) {
    case 2:
      return { rows: 2, cols: 2 };
    case 3:
      return { rows: 3, cols: 4 };
    case 4:
      return { rows: 4, cols: 4 };
    case 6:
      return { rows: 6, cols: 6 };
    default:
      return { rows: 4, cols: 4 };
  }
}

function createBoard(size) {
  currentBoardSize = size;
  const { rows, cols } = getBoardDimensions(size);
  const numCards = rows * cols;
  totalPairs = numCards / 2;

  moves = 0;
  matchedPairs = 0;
  flippedCards = [];
  movesSpan.textContent = `Moves: ${moves}`;
  gameBoard.innerHTML = "";
  winModal.style.display = "none";

  gameBoard.style.gridTemplateColumns = `repeat(${cols}, 1fr)`;
  gameBoard.style.gridTemplateRows = `repeat(${rows}, 1fr)`;

  const emojiSubset = shuffle(emojis).slice(0, totalPairs);
  const cardEmojis = shuffle([...emojiSubset, ...emojiSubset]);

  for (let i = 0; i < numCards; i++) {
    const card = document.createElement("div");
    card.classList.add("card");
    card.dataset.emoji = cardEmojis[i];

    const cardInner = document.createElement("div");
    cardInner.classList.add("card-inner");

    const cardFront = document.createElement("div");
    cardFront.classList.add("card-front");

    const cardBack = document.createElement("div");
    cardBack.classList.add("card-back");
    cardBack.textContent = cardEmojis[i];

    cardInner.appendChild(cardFront);
    cardInner.appendChild(cardBack);
    card.appendChild(cardInner);

    card.addEventListener("click", handleCardClick);
    gameBoard.appendChild(card);
  }
}

function handleCardClick(event) {
  const clickedCard = event.currentTarget;

  if (
    flippedCards.length < 2 &&
    !clickedCard.classList.contains("flipped") &&
    !clickedCard.classList.contains("matched")
  ) {
    clickedCard.classList.add("flipped");
    flippedCards.push(clickedCard);

    if (flippedCards.length === 2) {
      moves++;
      movesSpan.textContent = `Moves: ${moves}`;
      checkForMatch();
    }
  }
}

function checkForMatch() {
  const [card1, card2] = flippedCards;

  if (card1.dataset.emoji === card2.dataset.emoji) {
    card1.classList.add("matched");
    card2.classList.add("matched");
    matchedPairs++;
    flippedCards = [];
    if (matchedPairs === totalPairs) {
      setTimeout(showWinModal, 500);
    }
  } else {
    setTimeout(() => {
      card1.classList.remove("flipped");
      card2.classList.remove("flipped");
      flippedCards = [];
    }, 1000);
  }
}

function showWinModal() {
  finalMovesSpan.textContent = `You completed the game in ${moves} moves.`;
  winModal.style.display = "block";
}

function showMenu() {
  gameScreen.style.display = "none";
  winModal.style.display = "none";
  leaderboardScreen.style.display = "none";
  setupScreen.style.display = "block";
}

function showLeaderboard(size = 4) {
  // Default to 4x4
  setupScreen.style.display = "none";
  gameScreen.style.display = "none";
  winModal.style.display = "none";
  leaderboardScreen.style.display = "flex";
  displayLeaderboard(size);
}

function startGame(size) {
  setupScreen.style.display = "none";
  leaderboardScreen.style.display = "none";
  gameScreen.style.display = "flex";
  createBoard(size);
}

sizeButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const size = parseInt(button.dataset.size, 10);
    startGame(size);
  });
});

leaderboardForm.addEventListener("submit", (e) => {
  e.preventDefault();
  const initials = initialsInput.value.toUpperCase();
  if (initials) {
    addScore(currentBoardSize, initials, moves);
    initialsInput.value = "";
    showLeaderboard(currentBoardSize);
  }
});

leaderboardFilterButtons.forEach((button) => {
  button.addEventListener("click", () => {
    const size = parseInt(button.dataset.size, 10);
    displayLeaderboard(size);
  });
});

menuBtn.addEventListener("click", showMenu);
leaderboardBtn.addEventListener("click", () => showLeaderboard()); // Show with default size
backToMenuBtn.addEventListener("click", showMenu);

// Initial setup
showMenu();

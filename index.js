let sockAddr = 'ws://' + window.location.host + '/ws';
let socket = new WebSocket(sockAddr);

// Get canvas-related DOM variables.
let canvas = document.querySelector('canvas');
let ctx = canvas.getContext('2d');

// Draw black background.
ctx.fillStyle = "black";
ctx.fillRect(0, 0, canvas.width, canvas.height);
ctx.fill();

// Allow grabbing cursor and releasing it.
canvas.onclick = function() {
  canvas.requestPointerLock();
}

function sendMovement(e) {
  let a = Math.min(Math.abs(e.movementX) + 1, 127);
  let b = Math.min(Math.abs(e.movementY) + 1, 127);

  let c = [e.movementX > 0, e.movementY > 0].reduce((res, x) => res << 1 | x);
  socket.send(String.fromCharCode(65, a, b, c));
}

function sendMouseClick(e) {
  let msgType = e.type === 'mouseup' ? 68 : 67;
  socket.send(String.fromCharCode(msgType, e.button, 0, 0));
}

function sendMouseScroll(e) {
  let sign = (e.deltaY > 0) ? 0 : 2;
  socket.send(String.fromCharCode(69, 1, sign, 0));
}

document.addEventListener("pointerlockchange", () => {
  if (document.pointerLockElement === canvas) {
    document.addEventListener('mousemove', sendMovement);
    document.addEventListener('mouseup', sendMouseClick);
    document.addEventListener('mousedown', sendMouseClick);
    document.addEventListener('wheel', sendMouseScroll);
  } else {
    document.removeEventListener('mousemove', sendMovement);
    document.removeEventListener('mouseup', sendMouseClick);
    document.removeEventListener('mousedown', sendMouseClick);
    document.removeEventListener('wheel', sendMouseScroll);
  }
});

console.log("Attempting Connection...");

socket.onopen = () => {
  console.log("Successfully Connected");
};

socket.onclose = event => {
  console.log("Socket Closed Connection: ", event);
  socket.send("Client Closed!")
};

socket.onerror = error => {
  console.log("Socket Error: ", error);
};

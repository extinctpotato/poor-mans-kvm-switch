let sockAddr = 'ws://' + window.location.host + '/ws';
let socket = new WebSocket(sockAddr);

socket.binaryType = 'arraybuffer';

// Get canvas-related DOM variables.
let canvas = document.querySelector('canvas');
let ctx = canvas.getContext('2d');

// TODO: check if ContextMenu is correct.
const SPECIAL_KEY_MAPPING = {
  'ControlLeft': 0x80,
  'ShiftLeft': 0x81,
  'AltLeft': 0x82,
  'MetaLeft': 0x83,
  'ControlRight':  0x84,
  'ShiftRight': 0x85,
  'AltRight': 0x86,
  'MetaRight': 0x87,
  'ArrowUp': 0xDA,
  'ArrowDown': 0xD9,
  'ArrowLeft': 0xD8,
  'ArrowRight': 0xD7,
  'Backspace': 0xB2,
  'Tab': 0xB3,
  'ContextMenu': 0xED,
  'Escape': 0xB1,
  'PageUp': 0xD3,
  'PageDown': 0xD6,
  'Home': 0xD2,
  'End': 0xD5,
  'CapsLock': 0xC1,
  'PrintScreen': 0xCE,
  'ScrollLock': 0xCF,
  'Pause': 0xD0,
  'NumLock': 0xDB,
  'NumpadDivide': 0xDC,
  'NumpadMultiply': 0xDD,
  'NumpadSubtract': 0xDE,
  'NumpadAdd': 0xDF,
  'NumpadEnter': 0xE0,
  'Numpad0': 0xEA,
  'Numpad1': 0xE1,
  'Numpad2': 0xE2,
  'Numpad3': 0xE3,
  'Numpad4': 0xE4,
  'Numpad5': 0xE5,
  'Numpad6': 0xE6,
  'Numpad7': 0xE7,
  'Numpad8': 0xE8,
  'Numpad9': 0xE9,
  'NumpadComma': 0xEB,
  'F1': 0xC2,
  'F2': 0xC3,
  'F3': 0xC4,
  'F4': 0xC5,
  'F5': 0xC6,
  'F6': 0xC7,
  'F7': 0xC8,
  'F8': 0xC9,
  'F9': 0xCA,
  'F10': 0xCB,
  'F11': 0xCC,
  'F12': 0xCD,
  'F13': 0xF0,
  'F14': 0xF1,
  'F15': 0xF2,
  'F16': 0xF3,
  'F17': 0xF4,
  'F18': 0xF5,
  'F19': 0xF6,
  'F20': 0xF7,
  'F21': 0xF8,
  'F22': 0xF9,
  'F23': 0xFA,
  'F24': 0xFB,


};

// Draw black background.
ctx.fillStyle = "black";
ctx.fillRect(0, 0, canvas.width, canvas.height);
ctx.fill();

// Allow grabbing cursor and releasing it.
canvas.onclick = function() {
  canvas.requestPointerLock();
}

function makeMsg(msgType, ...args) {
  let msg = new Uint8Array(4);

  msg[0] = msgType;

  for (let i = 1; i < 4; i++) {
    msg[i] = args[i-1] || 0;
  }

  return msg;
}

function sendMovement(e) {
  let a = Math.min(Math.abs(e.movementX) + 1, 127);
  let b = Math.min(Math.abs(e.movementY) + 1, 127);

  let c = [e.movementX > 0, e.movementY > 0].reduce((res, x) => res << 1 | x);
  socket.send(makeMsg(65, a, b, c));
}

function sendMouseClick(e) {
  let msgType = e.type === 'mouseup' ? 68 : 67;
  socket.send(makeMsg(msgType, e.button));
}

function sendMouseScroll(e) {
  let sign = (e.deltaY > 0) ? 0 : 2;
  socket.send(makeMsg(69, 1, sign));
}

function sendKey(e) {
  let key;
  let msgType = e.type === 'keydown' ? 70 : 71;

  if (e.code in SPECIAL_KEY_MAPPING) {
    key = SPECIAL_KEY_MAPPING[e.code];
  } else if (e.key.length === 1) {
    key = e.key.charCodeAt(0);
  } else {
    console.log('Unsupported key type', key);
    return
  }

  console.log('Sending', key, e);

  socket.send(makeMsg(msgType, key));
}

document.addEventListener("pointerlockchange", () => {
  if (document.pointerLockElement === canvas) {
    document.addEventListener('mousemove', sendMovement);
    document.addEventListener('mouseup', sendMouseClick);
    document.addEventListener('mousedown', sendMouseClick);
    document.addEventListener('wheel', sendMouseScroll);
    document.addEventListener('keyup', sendKey);
    document.addEventListener('keydown', sendKey);
  } else {
    document.removeEventListener('mousemove', sendMovement);
    document.removeEventListener('mouseup', sendMouseClick);
    document.removeEventListener('mousedown', sendMouseClick);
    document.removeEventListener('wheel', sendMouseScroll);
    document.removeEventListener('keyup', sendKey);
    document.removeEventListener('keydown', sendKey);
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

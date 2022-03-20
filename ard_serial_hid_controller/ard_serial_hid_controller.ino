#include "Mouse.h"
#include "Keyboard.h"

byte buf[5];

struct IncomingMsg {
  byte command;
  byte arg1;
  byte arg2;
  byte arg3;
};

IncomingMsg incoming_msg;

char bitToMultiplier(byte b, int n) {
  if (b & (1 << n)) {
    return 1;
  } else {
    return -1;
  }
}

void deserialize_msg(byte* msg, struct IncomingMsg *i_msg) {
  i_msg->command = msg[0];
  i_msg->arg1 = msg[1];
  i_msg->arg2 = msg[2];
  i_msg->arg3 = msg[3];
}

void print_msg(struct IncomingMsg *i_msg) {
  Serial.println("=== BEGIN print_msg ===");
  Serial.println(i_msg->command, DEC);
  Serial.println(i_msg->arg1, DEC);
  Serial.println(i_msg->arg2, DEC);
  Serial.println(i_msg->arg3, DEC);
  Serial.println("=== END print_msg ===");
}

void handle_msg(struct IncomingMsg *i_msg) {
  char xMul, yMul;

  switch (i_msg->command) {
    case 10:
      break;
    case 65:
      // +--------+--------+-------------+
      // | X axis | Y axis | i_msg->arg3 |
      // +--------+--------+-------------+
      // | Neg    | Neg    |           0 |
      // | Neg    | Pos    |           1 |
      // | Pos    | Neg    |           2 |
      // | Pos    | Pos    |           3 |
      // +--------+--------+-------------+
      xMul = bitToMultiplier(i_msg->arg3, 1);
      yMul = bitToMultiplier(i_msg->arg3, 0);

      Mouse.move((char)(i_msg->arg1)*xMul, (char)(i_msg->arg2)*yMul, 0);
      break;
    case 66:
      switch(i_msg->arg1) {
        case 0:
          Mouse.click(MOUSE_LEFT);
          break;
        case 2:
          Mouse.click(MOUSE_RIGHT);
          break;
        case 1:
          Mouse.click(MOUSE_MIDDLE);
          break;
      }
      break;
    case 67:
      switch(i_msg->arg1) {
        case 0:
          Mouse.press(MOUSE_LEFT);
          break;
        case 2:
          Mouse.press(MOUSE_RIGHT);
          break;
        case 1:
          Mouse.press(MOUSE_MIDDLE);
          break;
      }
      break;
    case 68:
      switch(i_msg->arg1) {
        case 0:
          Mouse.release(MOUSE_LEFT);
          break;
        case 2:
          Mouse.release(MOUSE_RIGHT);
          break;
        case 1:
          Mouse.release(MOUSE_MIDDLE);
          break;
      }
      break;
    case 69:
      Mouse.move(0, 0, i_msg->arg1 * (i_msg->arg2 - 1));
      break;
    case 70:
      Keyboard.press(i_msg->arg1);
      break;
    case 71:
      Keyboard.release(i_msg->arg1);
      break;
    case 72:
      Keyboard.releaseAll();
      break;
  }
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  Serial1.begin(115200);
  Mouse.begin();
  Keyboard.begin();
}

void loop() {
  memset(&buf, 0, sizeof(buf));
  memset(&incoming_msg, 0, sizeof(incoming_msg));

  while (Serial1.available() == 0);
  Serial1.readBytes(buf, 4);

  deserialize_msg(buf, &incoming_msg);

  print_msg(&incoming_msg);
  handle_msg(&incoming_msg);
}

#include "Mouse.h"

char buf[5];

struct IncomingMsg {
  char command;
  char arg1;
  char arg2;
  char arg3;
};

IncomingMsg incoming_msg;

char bitToMultiplier(char b, int n) {
  if (b & (1 << n)) {
    return 1;
  } else {
    return -1;
  }
}

void deserialize_msg(char* msg, struct IncomingMsg *i_msg) {
  i_msg->command = msg[0];
  i_msg->arg1 = msg[1];
  i_msg->arg2 = msg[2];
  i_msg->arg3 = msg[3];
}

void serialize_msg(char* msg, struct IncomingMsg *i_msg) {
  snprintf(msg, 5, "%c%c%c%c",
      i_msg->command,
      i_msg->arg1,
      i_msg->arg2,
      i_msg->arg3
      );
}

void print_msg(struct IncomingMsg *i_msg) {
  Serial.println("=== BEGIN print_msg ===");
  Serial.println(i_msg->command, DEC);
  Serial.println(i_msg->arg1, DEC);
  Serial.println(i_msg->arg2, DEC);
  Serial.println(i_msg->arg3, DEC);
  Serial.println("=== END print_msg ===");
}

void print_msg_c(char* msg) {
  Serial.println("=== BEGIN print_msg_c ===");
  for (int i = 0; i < 4; i++) {
    Serial.print(i, DEC);
    Serial.print(": ");
    Serial.println(msg[i], DEC);
  }
  Serial.println("=== END print_msg_c ===");
}

void msg_demo() {
  struct IncomingMsg* i_msg = malloc(sizeof(struct IncomingMsg));
  struct IncomingMsg* i_msg_p = malloc(sizeof(struct IncomingMsg));

  char i_msg_c[5];

  i_msg->command = 2;
  i_msg->arg1 = 1;
  i_msg->arg2 = 3;
  i_msg->arg3 = 7;

  print_msg(i_msg);
  serialize_msg(i_msg_c, i_msg);
  Serial.println("serialized");

  print_msg_c(i_msg_c);

  deserialize_msg(i_msg_c, i_msg_p);
  print_msg(i_msg_p);

  free(i_msg);
  free(i_msg_p);
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

      Mouse.move((i_msg->arg1)*xMul, (i_msg->arg2)*yMul, 0);
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
  }
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(115200);
  Serial1.begin(115200);
  Mouse.begin();
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

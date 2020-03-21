PBDefine =  [[
syntax = "proto2";

message GameMessage {
  required string action = 1;
  optional bytes  data   = 2;
}

message RequestLoginMessage {
  required string token = 1;
}

message CommonMessage {

}

]]

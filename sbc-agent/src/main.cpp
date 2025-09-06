#include <iostream>
int camera_open();
int camera_read();
int webrtc_start();
int webrtc_send_frame();
int telem_open();
int telem_pump();

int main(){
  std::cout << "sbc-agent boot\n";
  camera_open();
  webrtc_start();
  telem_open();
  for(int i=0;i<3;++i){ camera_read(); webrtc_send_frame(); telem_pump(); }
  return 0;
}

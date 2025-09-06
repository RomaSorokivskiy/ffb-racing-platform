#include <iostream>
int init_video();
int pump_video();
namespace ffb { void init(); void send_torque(float); }
int net_connect();
int net_poll();

int main(){
  std::cout << "client boot\n";
  net_connect();
  init_video();
  ffb::init();
  // main loop (stub)
  for(int i=0;i<3;++i){ pump_video(); net_poll(); }
  return 0;
}

#include <iostream>
float compute_ffb(float steer_norm, float yaw_rate);
int main(){
  std::cout << "orchestrator boot\n";
  std::cout << "ffb sample: " << compute_ffb(0.2f, 10.0f) << "\n";
  return 0;
}

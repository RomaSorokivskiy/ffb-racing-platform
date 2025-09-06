#include <iostream>
float compute_ffb(float steer_norm, float yaw_rate){
  // простий демо-FFB: пружина + демпфер
  const float k_spring = 2.0f;       // Нм/од.
  const float k_damper = 0.05f;      // Нм/(deg/s)
  return (-k_spring*steer_norm) - (k_damper*yaw_rate);
}

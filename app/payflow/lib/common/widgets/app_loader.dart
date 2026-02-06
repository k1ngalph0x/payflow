import 'package:flutter/material.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

class AppLoader extends StatelessWidget {
  final Color color;
  final double size;
  const AppLoader({super.key, required this.color, required this.size});

  @override
  Widget build(BuildContext context) {
    return LoadingAnimationWidget.fourRotatingDots(color: color, size: size);
  }
}

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:payflow/utils/app_routes.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    navigate();
  }

  Future<void> navigate() async {
    await Future.delayed(const Duration(seconds: 2));

    Get.toNamed(AppRoutes.home);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold();
  }
}

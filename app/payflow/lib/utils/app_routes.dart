import 'package:get/get.dart';
import '../common/views/splash_screen.dart';
import 'package:payflow/common/views/home_screen.dart';

class AppRoutes {
  static const String home = "/home";
  static const String splash = "/splash";

  static final routes = [
    GetPage(
      name: splash,
      page: () => const SplashScreen(),
      transition: Transition.fade,
    ),

    GetPage(
      name: home,
      page: () => const HomeScreen(),
      transition: Transition.fade,
    ),
  ];
}

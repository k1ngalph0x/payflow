import 'package:get/get.dart';
import 'package:payflow/common/views/role_selection_screen.dart';
import 'package:payflow/common/views/signin_screen.dart';
import 'package:payflow/common/views/signup_screen.dart';
import 'package:payflow/common/views/splash_screen.dart';
import 'package:payflow/merchant/views/merchant_home_screen.dart';
import 'package:payflow/merchant/views/merchant_onboarding_screen.dart';
import 'package:payflow/user/views/user_home_screen.dart';

class AppRoutes {
  static const String splash = '/splash';
  static const String roleSelection = '/role-selection';
  static const String signIn = '/signin';
  static const String signUp = '/signup';
  static const String userHome = '/user-home';
  static const String merchantHome = '/merchant-home';
  static const String merchantOnboarding = '/merchant-onboarding';

  static List<GetPage> routes = [
    GetPage(name: splash, page: () => const SplashScreen()),
    GetPage(name: roleSelection, page: () => const RoleSelectionScreen()),
    GetPage(name: signIn, page: () => const SignInScreen()),
    GetPage(name: signUp, page: () => const SignUpScreen()),
    GetPage(name: userHome, page: () => const UserHomeScreen()),
    GetPage(name: merchantHome, page: () => const MerchantHomeScreen()),
    GetPage(
      name: merchantOnboarding,
      page: () => const MerchantOnboardingScreen(),
    ),
  ];
}

import 'package:flutter/cupertino.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPreferenceUtil {
  static Future<SharedPreferences> get instance async =>
      prefsInstance ??= await SharedPreferences.getInstance();
  static SharedPreferences? prefsInstance;

  static Future<SharedPreferences> init() async {
    SharedPreferences prefsInstance = await instance;
    return prefsInstance;
  }

  static Future<void> remove(String key) async {
    SharedPreferences prefs = await instance;
    await prefs.remove(key);
    debugPrint("Removed key: $key");
  }

  static Future<void> prefranceClear() async {
    SharedPreferences prefs = await instance;
    prefs.clear();
    debugPrint("All Preferences are Cleared.");
  }

  static Future<void> setString(String key, String value) async {
    SharedPreferences prefs = await instance;
    await prefs.setString(key, value);
  }

  static Future<void> setBool(String key, bool value) async {
    SharedPreferences prefs = await instance;
    await prefs.setBool(key, value);
  }

  static Future<void> setInt(String key, int value) async {
    SharedPreferences prefs = await instance;
    await prefs.setInt(key, value);
  }

  static String getString(String key, [String? defValue]) {
    return prefsInstance?.getString(key) ?? defValue ?? "";
  }

  static bool getBool(String key, [bool? defValue]) {
    return prefsInstance?.getBool(key) ?? defValue ?? false;
  }

  static int getInt(String key, [int? defValue]) {
    return prefsInstance?.getInt(key) ?? defValue ?? 0;
  }
}

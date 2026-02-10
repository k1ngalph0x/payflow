import 'package:get/get.dart';
import 'package:payflow/api/api_repo.dart';
import 'package:payflow/common/models/merchant_list_model.dart';

class MerchantController extends GetxController {
  final _apiRepo = ApiRepo();

  final RxBool isLoading = false.obs;
  final RxList<MerchantListModel> merchants = <MerchantListModel>[].obs;
  final Rxn<MerchantListModel> selectedMerchant = Rxn<MerchantListModel>();

  @override
  void onInit() {
    super.onInit();
    fetchMerchants();
  }

  Future<void> fetchMerchants() async {
    try {
      isLoading.value = true;

      final response = await _apiRepo.getMerchantList();

      if (response['merchants'] != null) {
        final List<dynamic> merchantList = response['merchants'];
        merchants.value = merchantList
            .map((json) => MerchantListModel.fromJson(json))
            .toList();
      }
    } catch (e) {
      print('Fetch merchants error: $e');
    } finally {
      isLoading.value = false;
    }
  }

  void selectMerchant(MerchantListModel merchant) {
    selectedMerchant.value = merchant;
  }
}

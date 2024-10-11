import 'package:get_it/get_it.dart';
import 'package:image_flipper_gui/app/services/image_service.dart';
import 'package:image_flipper_gui/domain/interfaces/i_image_service.dart';

final getIt = GetIt.instance;

void setupDi() {
  getIt.registerSingleton<IImageService>(ImageService());
}

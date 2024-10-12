import 'package:file_picker/file_picker.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_state.dart';
import 'package:image_flipper_gui/domain/entities/flip_action.dart';
import 'package:image_flipper_gui/domain/interfaces/i_image_service.dart';

class FlipperCubit extends Cubit<FlipperState> {
  final IImageService _imageService;

  FlipperCubit(this._imageService) : super(FlipperInitial());

  Future<void> loadFolder() async {
    try {
      emit(FlipperLoadFolderInProgress());

      final selectedFolder = await _pickFolder();

      emit(
        FlipperLoadFolderSuccess(
          images: _imageService.filterImages(selectedFolder),
        ),
      );
    } catch (e) {
      emit(FlipperError(message: e.toString()));
    }
  }

  Future<void> previewImagesFlip(FlipAction flip) async {
    return switch (state) {
      FlipperLoadFolderSuccess(:final images) =>
        emit(FlipperPreviewFlipImagesSuccess(
          images: images,
          action: flip,
        )),
      FlipperPreviewFlipImagesSuccess(:final images) =>
        emit(FlipperPreviewFlipImagesSuccess(
          images: images,
          action: flip,
        )),
      FlipperSaveImagesSuccess(:final images) =>
        emit(FlipperPreviewFlipImagesSuccess(
          images: images,
          action: flip,
        )),
      _ => emit(FlipperError(message: 'No images to preview'))
    };
  }

  Future<void> saveFlippedImages() async {
    try {
      if (state.noFlipActionSelected) {
        emit(FlipperError(message: 'No flip action selected'));
        return;
      }

      final selectedFolder = await _pickFolder();

      switch (state) {
        case FlipperPreviewFlipImagesSuccess(:final images, :final action):
          {
            emit(FlipperSaveImagesInProgress());
            await _imageService.flipImages(
              images,
              action,
              selectedFolder,
            );
            emit(FlipperSaveImagesSuccess(
              images: images,
              action: action,
              outputDir: selectedFolder,
            ));
            break;
          }
        default:
          {
            emit(FlipperError(message: 'No images to save'));
          }
      }
    } catch (e) {
      emit(FlipperError(message: e.toString()));
    }
  }

  void discardSelectedImages() {
    emit(FlipperInitial());
  }

  Future<String> _pickFolder() async {
    final selectedFolder = await FilePicker.platform.getDirectoryPath(
      dialogTitle: 'Select a folder',
    );

    if (selectedFolder == null) {
      emit(FlipperError(message: 'No folder selected'));
      return '';
    }

    return selectedFolder;
  }
}

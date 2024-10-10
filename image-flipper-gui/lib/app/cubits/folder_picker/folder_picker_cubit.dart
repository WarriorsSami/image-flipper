import 'package:file_picker/file_picker.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/folder_picker/folder_picker_state.dart';
import 'package:image_flipper_gui/app/helper.dart';

class FolderCubit extends Cubit<FolderState> {
  FolderCubit() : super(FolderInitial());

  Future<void> loadFolder() async {
    try {
      emit(FolderLoading());

      final selectedFolder = await FilePicker.platform.getDirectoryPath(
        dialogTitle: 'Select a folder',
      );

      if (selectedFolder == null) {
        emit(FolderError(message: 'No folder selected'));
        return;
      }

      emit(FolderLoaded(images: filterImages(selectedFolder)));
    } catch (e) {
      emit(FolderError(message: e.toString()));
    }
  }
}

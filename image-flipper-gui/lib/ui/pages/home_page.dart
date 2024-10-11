import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_cubit.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_state.dart';
import 'package:image_flipper_gui/domain/entities/flip_action.dart';
import 'package:image_flipper_gui/ui/widgets/images_widget.dart';

class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Image flipper'),
        centerTitle: true,
        backgroundColor: Theme.of(context).colorScheme.tertiary,
        foregroundColor: Colors.white,
      ),
      body: BlocBuilder<FlipperCubit, FlipperState>(
        builder: (context, state) {
          return Center(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: <Widget>[
                const SizedBox(height: 10),
                Text(
                  switch (context.read<FlipperCubit>().state) {
                    FlipperSaveImagesSuccess(:final outputDir) =>
                      'Images saved successfully to $outputDir!',
                    FlipperSaveImagesInProgress() =>
                      'Saving images... Please wait!',
                    _ => 'Here you can flip your images!',
                  },
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
                const SizedBox(height: 30),
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    ElevatedButton(
                      onPressed: () async {
                        await context.read<FlipperCubit>().loadFolder();
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Theme.of(context).colorScheme.primary,
                        foregroundColor: Colors.white,
                      ),
                      child: const Text('Select source folder'),
                    ),
                    const SizedBox(width: 10),
                    ElevatedButton(
                      onPressed: context
                              .read<FlipperCubit>()
                              .state
                              .noFlipActionSelected
                          ? null
                          : () async {
                              await context
                                  .read<FlipperCubit>()
                                  .saveFlippedImages();
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Theme.of(context).colorScheme.primary,
                        foregroundColor: Colors.white,
                      ),
                      child: const Text('Save to folder'),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                ElevatedButton(
                  onPressed: () {
                    context.read<FlipperCubit>().discardSelectedImages();
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Theme.of(context).colorScheme.primary,
                    foregroundColor: Colors.white,
                  ),
                  child: const Text('Discard images'),
                ),
                const SizedBox(height: 20),
                SizedBox(
                  width: 200,
                  child: Column(
                    children: [
                      RadioListTile<FlipAction>(
                        title: const Text('Horizontally'),
                        value: FlipAction.horizontal,
                        groupValue: switch (
                            context.read<FlipperCubit>().state) {
                          FlipperPreviewFlipImagesSuccess(:final action) ||
                          FlipperSaveImagesSuccess(:final action) =>
                            action,
                          _ => null,
                        },
                        onChanged:
                            context.read<FlipperCubit>().state.noFolderSelected
                                ? null
                                : (FlipAction? value) {
                                    context
                                        .read<FlipperCubit>()
                                        .previewImagesFlip(value!);
                                  },
                      ),
                      RadioListTile<FlipAction>(
                        title: const Text('Vertically'),
                        value: FlipAction.vertical,
                        groupValue: switch (
                            context.read<FlipperCubit>().state) {
                          FlipperPreviewFlipImagesSuccess(:final action) ||
                          FlipperSaveImagesSuccess(:final action) =>
                            action,
                          _ => null,
                        },
                        onChanged:
                            context.read<FlipperCubit>().state.noFolderSelected
                                ? null
                                : (FlipAction? value) {
                                    context
                                        .read<FlipperCubit>()
                                        .previewImagesFlip(value!);
                                  },
                      ),
                      RadioListTile<FlipAction>(
                        title: const Text('Both'),
                        value: FlipAction.both,
                        groupValue: switch (
                            context.read<FlipperCubit>().state) {
                          FlipperPreviewFlipImagesSuccess(:final action) ||
                          FlipperSaveImagesSuccess(:final action) =>
                            action,
                          _ => null,
                        },
                        onChanged:
                            context.read<FlipperCubit>().state.noFolderSelected
                                ? null
                                : (FlipAction? value) {
                                    context
                                        .read<FlipperCubit>()
                                        .previewImagesFlip(value!);
                                  },
                      ),
                      RadioListTile<FlipAction>(
                        title: const Text('Original'),
                        value: FlipAction.original,
                        groupValue: switch (
                            context.read<FlipperCubit>().state) {
                          FlipperPreviewFlipImagesSuccess(:final action) ||
                          FlipperSaveImagesSuccess(:final action) =>
                            action,
                          _ => null,
                        },
                        onChanged:
                            context.read<FlipperCubit>().state.noFolderSelected
                                ? null
                                : (FlipAction? value) {
                                    context
                                        .read<FlipperCubit>()
                                        .previewImagesFlip(value!);
                                  },
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 20),
                const Expanded(
                  child: ImagesWidget(),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}

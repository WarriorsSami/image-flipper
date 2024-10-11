import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:image_flipper_gui/app/cubits/flipper/flipper_cubit.dart';
import 'package:image_flipper_gui/di.dart';
import 'package:image_flipper_gui/domain/interfaces/i_image_service.dart';
import 'package:image_flipper_gui/ui/pages/home_page.dart';

class ImageFlipperApp extends StatelessWidget {
  const ImageFlipperApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Image Flipper',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        primarySwatch: Colors.teal,
        colorScheme: ColorScheme.fromSeed(
          seedColor: Colors.lightGreen,
          secondary: Colors.lightBlue,
        ),
        useMaterial3: true,
      ),
      home: BlocProvider<FlipperCubit>(
        create: (context) => FlipperCubit(getIt<IImageService>()),
        child: const HomePage(),
      ),
    );
  }
}

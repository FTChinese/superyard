package apps

import (
	"testing"
)

const mockGradleFile = `
apply plugin: 'com.android.application'
apply plugin: 'kotlin-android'
apply plugin: 'kotlin-android-extensions'
apply plugin: 'kotlin-kapt'
apply plugin: 'io.fabric'

// See:
// https://stackoverflow.com/questions/37101589/how-to-read-a-properties-files-and-use-the-values-in-project-gradle-script/37101792
// https://stackoverflow.com/questions/45586839/android-read-build-gradle-properties-inside-class
// https://medium.com/@abhi007tyagi/storing-api-keys-using-android-ndk-6abb0adcadad
def props = new Properties()
file("$rootDir/config.properties").withInputStream { props.load(it) }

// Creates a variable called keystorePropertiesFile, and initializes it to the
// keystore.properties file.
def keystorePropertiesFile = rootProject.file("keystore.properties")

// Initializes a new Properties() object called keystoreProperties.
def keystoreProperties = new Properties()

// Loads the keystore.properties file into the keystoreProperties object.
keystoreProperties.load(new FileInputStream(keystorePropertiesFile))

androidExtensions {
    experimental = true
}

android {
    compileSdkVersion 29
    compileOptions {
        sourceCompatibility 1.8
        targetCompatibility 1.8
    }
    kotlinOptions {
        jvmTarget = "1.8"
    }
    defaultConfig {
        applicationId "com.ft.ftchinese"
        minSdkVersion 21
        targetSdkVersion 29
        versionCode 34
        versionName "3.2.9"

        setProperty("archivesBaseName", "ftchinese-v$versionName")
        testInstrumentationRunner "androidx.test.runner.AndroidJUnitRunner"
        buildConfigField "String", "WX_SUBS_APPID", props.getProperty("wechat.subs.appId")

        buildConfigField "String", "BASE_URL_STANDARD", props.getProperty("base_url.standard")
        buildConfigField "String", "BASE_URL_PREMIUM", props.getProperty("base_url.premium")
        buildConfigField "String", "BASE_URL_B2B", props.getProperty("base_url.b2b")
        externalNativeBuild {
            cmake {
                cppFlags "-std=c++11"
            }
        }
        // Export database schema to app/schemas directory.
        javaCompileOptions {
            annotationProcessorOptions {
                arguments = ["room.schemaLocation": "$projectDir/schemas".toString()]
            }
        }
    }
    dataBinding {
        enabled true
    }
    signingConfigs {
        release {
            keyAlias keystoreProperties['keyAlias']
            keyPassword keystoreProperties['keyPassword']
            storeFile file("$rootDir/android.jks")
            storePassword keystoreProperties['storePassword']
        }
    }
    buildTypes {
        release {
            // Adds the "release" signing configuration to the release build type.
            signingConfig signingConfigs.release
            minifyEnabled false
            proguardFiles getDefaultProguardFile('proguard-android.txt'), 'proguard-rules.pro'
            debuggable false
            buildConfigField "String", "STRIPE_KEY", props.getProperty("stripe.live")
            buildConfigField "String", "ACCESS_TOKEN", props.getProperty("access_token.live")
        }
        debug {
            debuggable true
            buildConfigField "String", "STRIPE_KEY", props.getProperty("stripe.test")
            buildConfigField "String", "ACCESS_TOKEN", props.getProperty("access_token.test")
        }
    }

    flavorDimensions "appStore"
    productFlavors {
        // For google play. Take this as the official version.
        play {
            dimension "appStore"
        }
        huawei {
            dimension "appStore"
            versionNameSuffix "-huawei"
        }
        sanliuling {
            dimension "appStore"
            versionNameSuffix "-360"
        }
        ftc {
            dimension "appStore"
            versionNameSuffix "-ftc"
        }
        samsung {
            dimension "appStore"
            versionNameSuffix "-samsung"
        }
        anzhi {
            dimension "appStore"
            versionNameSuffix "-anzhi"
        }
        standard {
            dimension "appStore"
            versionNameSuffix "-standard"
        }
        premium {
            dimension "appStore"
            versionNameSuffix "-premium"
        }
        b2b {
            dimension "appStore"
            versionNameSuffix "-b2b"
        }
    }

    testOptions {
        unitTests {
            returnDefaultValues = true
            includeAndroidResources = true
        }
    }

//    externalNativeBuild {
//        cmake {
//            path "CMakeLists.txt"
//        }
//    }
}

//kapt {
//    generateStubs = true
//}

dependencies {
    def room_version = '2.2.1'
    def lifecycle_version = "2.2.0-rc03"
//    def nav_version = "2.1.0-alpha05"
    def nav_version_ktx = "2.2.0-rc04"
    def exo_player = "2.11.1"

    implementation "org.jetbrains.kotlin:kotlin-stdlib:$kotlin_version"
    implementation "org.jetbrains.anko:anko:$anko_version"
    implementation "org.jetbrains.kotlinx:kotlinx-coroutines-android:$coroutines_version"

    // Firebase
    implementation 'com.google.firebase:firebase-core:17.2.2'
    implementation 'com.google.firebase:firebase-analytics:17.2.2'
    implementation 'com.google.firebase:firebase-iid:20.0.2'
    implementation 'com.google.firebase:firebase-messaging:20.1.0'
    implementation 'com.google.android.gms:play-services-analytics:17.0.0'
    implementation 'com.crashlytics.sdk.android:crashlytics:2.10.1'

    // Payment
    implementation(name: 'alipaySdk-15.6.5-20190718211148', ext: 'aar')
    implementation 'com.stripe:stripe-android:10.1.0'


    // ViewModel and LiveData
    implementation "androidx.lifecycle:lifecycle-extensions:$lifecycle_version"

    implementation "androidx.lifecycle:lifecycle-viewmodel-ktx:$lifecycle_version"
    implementation "androidx.lifecycle:lifecycle-livedata-ktx:$lifecycle_version"
    implementation "androidx.fragment:fragment-ktx:1.2.0-rc05"
//    kapt "androidx.lifecycle:lifecycle-compiler:$lifecycle_version"
    implementation "androidx.lifecycle:lifecycle-common-java8:$lifecycle_version"

    implementation 'android.arch.work:work-runtime:1.0.1'
    implementation 'androidx.core:core-ktx:1.1.0'
    implementation 'androidx.legacy:legacy-support-v4:1.0.0'

    //  UI
    implementation "androidx.appcompat:appcompat:1.1.0"
    implementation "com.google.android.material:material:1.2.0-alpha03"
    implementation 'androidx.constraintlayout:constraintlayout:1.1.3'
    implementation "androidx.browser:browser:1.2.0"
    implementation "androidx.preference:preference:1.1.0"
    implementation "androidx.legacy:legacy-support-v4:$support_version"
    implementation "androidx.recyclerview:recyclerview:1.1.0"
//    implementation "androidx.coordinatorlayout:coordinatorlayout:1.1.0"
    implementation "androidx.cardview:cardview:1.0.0"
    implementation "com.makeramen:roundedimageview:2.3.0"
    implementation "androidx.navigation:navigation-fragment-ktx:$nav_version_ktx"
    implementation "androidx.navigation:navigation-ui-ktx:$nav_version_ktx"

    // ExoPlayer
    implementation "com.google.android.exoplayer:exoplayer:$exo_player"
    implementation "com.google.android.exoplayer:extension-mediasession:$exo_player"

    //  Network
    implementation 'com.squareup.okhttp3:okhttp:4.2.2'
    implementation 'com.beust:klaxon:5.0.2'

    //  ORM
    implementation "androidx.room:room-runtime:$room_version"
    kapt "androidx.room:room-compiler:$room_version"
    implementation "androidx.room:room-coroutines:2.1.0-alpha04"

    // Wechat
    implementation 'com.tencent.mm.opensdk:wechat-sdk-android-with-mta:5.4.0'

    // Utilities
    implementation 'com.jakewharton.byteunits:byteunits:0.9.1'
    implementation 'org.threeten:threetenbp:1.4.0'
    implementation 'org.apache.commons:commons-math3:3.6.1'


    //    Test
    testImplementation 'junit:junit:4.12'
    testImplementation 'com.github.javafaker:javafaker:1.0.1'
    androidTestImplementation 'androidx.test:runner:1.2.0'
    androidTestImplementation 'androidx.test.espresso:espresso-core:3.2.0'
    testImplementation 'org.hamcrest:hamcrest-junit:2.0.0.0'
    testImplementation 'org.mockito:mockito-core:3.1.0'
    testImplementation 'org.robolectric:robolectric:4.3.1'
    testImplementation "androidx.room:room-testing:$room_version"
    // optional - Test helpers for LiveData
    testImplementation "androidx.arch.core:core-testing:2.1.0"
}

apply plugin: 'com.google.gms.google-services'
`

func TestExtractVersionCode(t *testing.T) {
	versionCode := ExtractVersionCode(mockGradleFile)

	t.Log(versionCode)
}

func TestGHRepo_GradleFile(t *testing.T) {
	repo := NewGHRepo()

	content, rErr := repo.GradleFile("v3.2.9")

	if rErr != nil {
		t.Error(rErr)
	}

	buildGradle, err := content.GetContent()

	if err != nil {
		t.Error(err)
	}

	t.Log(ExtractVersionCode(buildGradle))
}

func TestGHRepo_LatestRelease(t *testing.T) {
	repo := NewGHRepo()

	r, err := repo.LatestRelease()

	if err != nil {
		t.Error(err)
	}

	t.Log(r)
}

func TestGHRepo_SingleRelease(t *testing.T) {
	repo := NewGHRepo()

	r, err := repo.SingleRelease("v3.2.9")

	if err != nil {
		t.Error(err)
	}

	t.Log(r)
}

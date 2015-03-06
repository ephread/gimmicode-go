module.exports = function(grunt) {

  var sassFiles = {
    'public/stylesheets/application.css': 'assets/stylesheets/application.scss'
  }

  var uglifyFiles = {
    'public/javascripts/application.js': [
      'assets/javascripts/jquery-2.1.3.js',
      'assets/javascripts/gimmicode.js'
    ]
  }

  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    sass: {
      development: {
        options: {
          style: 'compressed',
          sourcemap: 'auto'
        },
        files: sassFiles
      },

      production: {
        options: {
          style: 'compressed',
          sourcemap: 'none'
        },
        files: sassFiles
      }
    },
    uglify: {
      development: {
        options: {
          sourceMap: true,
          compress: true
        },
        files: uglifyFiles
      },

      production: {
        options: {
          sourceMap: false,
          compress: true
        },
        files: uglifyFiles
      }
    },
    watch: {
      css: {
        files: '**/*.scss',
        tasks: ['sass:development']
      },
      js: {
        files: '**/*.js',
        tasks: ['uglify:development']
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-sass');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-uglify');

  grunt.registerTask('default', function () {
    grunt.task.run(['watch']);
  });

  grunt.registerTask('development', function () {
    grunt.task.run(['sass:development', 'uglify:development']);
  });

  grunt.registerTask('production', function () {
    grunt.task.run(['sass:production', 'uglify:production']);
  });
}

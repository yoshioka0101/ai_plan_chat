env "dev" {
  src = "file://."
  dev = "docker://mysql/8/dev"
  url = "mysql://root:password@localhost:3306/ai_chat_task"
  
  migration {
    dir = "file://."
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

